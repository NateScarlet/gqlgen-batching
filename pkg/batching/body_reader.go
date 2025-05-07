package batching

import (
	"bufio"
	"io"
)

func newBodyReader(
	r io.ReadCloser,
) *bodyReader {
	return &bodyReader{
		r,
		bufio.NewReader(r),
	}
}

type bodyReader struct {
	raw io.ReadCloser
	r   *bufio.Reader
}

// Close implements io.ReadCloser.
func (b *bodyReader) Close() error {
	return b.raw.Close()
}

// Read implements io.ReadCloser.
func (b *bodyReader) Read(p []byte) (n int, err error) {
	return b.r.Read(p)
}

func (b *bodyReader) IsArray() bool {
	var v, _ = b.r.Peek(1)
	return string(v) == "["
}
