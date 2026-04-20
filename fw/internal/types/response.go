package types

import (
	"io"
	"net/http"

	"github.com/Deimvis/go-ext/go1.25/xhttp"
)

// TODO: somehow link with request (in order to allow validation request-response pairing)
// e.g. add PairedResponse[T fw.Request] with RequestPair() T
// e.g. add RequestKey interface that implements Method() string + Path() string,
// and PairedResponse with RequestPair() RequestKey, but issue with URI params (dynamic parts of Path()), maybe PathWildcard() string, maybe RequestMatcher(fw.Request) instead of RequestKey
// Maybe both options as StronglyPairedResponse and PairedResponse
type Response interface {
	Code() int
	Header() xhttp.ConstHeader
	// BodyStream considers body being immutable
	// and returns always the same instance of body reader
	BodyStream() io.Reader
}

type BufferedResponse interface {
	Response
	BodyRaw() []byte
}

type ResponseReflect interface {
	Response
	SetCode(int) error
	SetHeader(http.Header) error
	SetBodyStream(io.Reader) error
}
