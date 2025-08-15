package batching

import (
	"encoding/json"
	"mime"
	"net/http"
	"slices"
	"sync"

	"github.com/NateScarlet/gqlgen-batching/internal/iterator"

	"github.com/99designs/gqlgen/graphql"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

var _ graphql.Transport = POST{}

// POST implements https://github.com/graphql/graphql-over-http/blob/main/rfcs/Batching.md
type POST struct {
	// Map of all headers that are added to graphql response. If not
	// set, only one header: Content-Type: application/json will be set.
	ResponseHeaders           map[string][]string
	ConcurrentLimitPerRequest int
}

// Supports implements graphql.Transport.
func (h POST) Supports(r *http.Request) bool {
	if r.Header.Get("Upgrade") != "" {
		return false
	}

	mediaType, _, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
	if err != nil {
		return false
	}

	if !(r.Method == "POST" && mediaType == "application/json") {
		return false
	}

	var body = newBodyReader(r.Body)
	r.Body = body
	return body.IsArray()
}

func (h POST) Do(w http.ResponseWriter, r *http.Request, exec graphql.GraphExecutor) {
	start := graphql.Now()
	ctx := r.Context()
	writeHeaders(w, h.ResponseHeaders)
	var paramsCollection []*graphql.RawParams
	err := readJSON(r.Body, &paramsCollection)
	r.Body.Close()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		gqlErr := gqlerror.Errorf(
			"json request body could not be decoded: %+v",
			err,
		)
		resp := exec.DispatchError(graphql.WithOperationContext(ctx, &graphql.OperationContext{}), gqlerror.List{gqlErr})
		writeJSON(w, resp)
		return
	}

	var encoder *json.Encoder
	var writeHeaderOnce sync.Once
	for response := range iterator.Parallel(ctx, h.ConcurrentLimitPerRequest, slices.Values(paramsCollection), func(params *graphql.RawParams) *graphql.Response {
		params.Headers = r.Header
		params.ReadTime = graphql.TraceTiming{
			Start: start,
			End:   graphql.Now(),
		}
		rc, OpErr := exec.CreateOperationContext(ctx, params)
		if OpErr != nil {
			writeHeaderOnce.Do(func() {
				// 只有发送第一个响应前的错误影响状态码
				w.WriteHeader(statusFor(OpErr))
			})
			return exec.DispatchError(graphql.WithOperationContext(ctx, rc), OpErr)
		}
		responseHandler, ctx := exec.DispatchOperation(ctx, rc)
		return responseHandler(ctx)
	}) {
		writeHeaderOnce.Do(func() {}) // 将要写入响应体了，无法再修改 header
		if encoder == nil {
			encoder = json.NewEncoder(w)
			w.Write([]byte("["))
			defer w.Write([]byte("]"))
		} else {
			w.Write([]byte(","))
		}
		var err = encoder.Encode(response)
		if err != nil {
			return
		}
	}
}
