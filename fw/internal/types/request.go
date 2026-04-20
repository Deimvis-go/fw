package types

import (
	"io"
	"net/http"

	"github.com/Deimvis/go-ext/go1.25/xhttp"
)

// TODO: rename *Reflect to Mutable*

// TODO: how to specify mime type if not Header?
// TODO: maybe header modifiers? like JSONHeaders (Content-Type, Accept for req if resp has json type) SetCookie headers
// they will be set on construction, but will be overridden with SetHeader()

type Request interface {
	Method() string
	Path() string
	QueryString() string
	Header() xhttp.ConstHeader
	// BodyStream considers body being immutable
	// and returns always the same instance of body reader
	BodyStream() io.Reader
}

type BufferedRequest interface {
	Request
	BodyRaw() []byte
	// TODO: maybe add NewBodyStream() io.Reader for clarity (<=> bytes.NewReader(BodyRaw()))
}

type RequestReflect interface {
	Request
	SetMethod(string) error
	SetPath(string) error
	// TODO: SetQueryString(string) error
	SetHeader(http.Header) error
	SetBodyStream(io.Reader) error
}
