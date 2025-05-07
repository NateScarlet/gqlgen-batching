package batching

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/errcode"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

func writeJSON[T *graphql.Response | []*graphql.Response](w io.Writer, response T) {
	var encoder = json.NewEncoder(w)
	var err = encoder.Encode(response)
	if err != nil {
		panic(err)
	}
}

func readJSON(r io.Reader, val interface{}) error {
	dec := json.NewDecoder(r)
	dec.UseNumber()
	return dec.Decode(val)
}

func statusFor(errs gqlerror.List) int {
	switch errcode.GetErrorKind(errs) {
	case errcode.KindProtocol:
		return http.StatusUnprocessableEntity
	default:
		return http.StatusOK
	}
}
