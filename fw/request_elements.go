package fw

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/Deimvis-go/fw/fw/fwheader"
	"github.com/Deimvis-go/fw/fw/internal/utils"
	"github.com/Deimvis/go-ext/go1.25/xhttp"
)

// Method

type RequestGET struct {
	requestMethodBound
}

func (r RequestGET) Method() string {
	return "GET"
}

type RequestPOST struct {
	requestMethodBound
}

func (r RequestPOST) Method() string {
	return "POST"
}

type RequestPUT struct {
	requestMethodBound
}

func (r RequestPUT) Method() string {
	return "PUT"
}

type RequestPATCH struct {
	requestMethodBound
}

func (r RequestPATCH) Method() string {
	return "PATCH"
}

type RequestDELETE struct {
	requestMethodBound
}

func (r RequestDELETE) Method() string {
	return "DELETE"
}

type RequestHEAD struct {
	requestMethodBound
}

func (r RequestHEAD) Method() string {
	return "HEAD"
}

// Path

type RequestPathBound struct{}

func (r RequestPathBound) SetPath(string) error {
	return errBoundRequestPath
}

var errBoundRequestPath = errors.New("request path is bound to request type and can't be set")

// Header

type RequestNoHeader struct{}

func (r RequestNoHeader) Header() xhttp.ConstHeader {
	// not returning nil, but const header wrap,
	// because it's safe to call methods over it.
	return xhttp.AsConstHeader(nil)
}

type RequestHeader[PresetT HeaderPreset] struct {
	// TODO: maybe make Headers unexported
	Headers http.Header `json:"-"`
	hp      PresetT
}

func (r *RequestHeader[PresetT]) Header() xhttp.ConstHeader {
	if len(r.Headers) == 0 {
		r.Headers = r.hp.New()
	}
	return xhttp.AsConstHeader(r.Headers)
}

func (r *RequestHeader[PresetT]) SetHeader(h http.Header) error {
	r.Headers = h
	return nil
}

// GOEXPERIMENT=aliastypeparams : make alias
// type RequestStructHeader[PresetT HeaderPreset, T any] = fwheader.Structured[PresetT, T]
type RequestStructHeader[PresetT HeaderPreset, T any] struct {
	fwheader.Structured[PresetT, T]
}

// URI

type RequestNoURI struct{}

type RequestURI[UriT any] struct {
	URI UriT
}

// Query

type RequestNoQuery struct{}

func (rq *RequestNoQuery) QueryString() string {
	return ""
}

type RequestQuery[QueryT any] struct {
	Query QueryT
}

func (rq *RequestQuery[QueryT]) GetQuery() *QueryT {
	return &rq.Query
}

func (rq *RequestQuery[QueryT]) QueryString() string {
	return utils.MakeQueryString(rq.Query)
}

// Body

type RequestNoBody struct{}

func (rb *RequestNoBody) BodyStream() io.Reader {
	return bytes.NewReader(nil)
}

func (rb *RequestNoBody) BodyRaw() []byte {
	return nil
}

func (rb *RequestNoBody) SetBodyStream(io.Reader) error {
	// ignore extra fields
	return nil
}

func (rb *RequestNoBody) SetBodyRaw([]byte) error {
	// ignore extra fields
	return nil
}

type RequestBodyJSON[JsonBodyT any] struct {
	Body JsonBodyT
	r    io.Reader
}

func (rb *RequestBodyJSON[BodyT]) BodyStream() io.Reader {
	if rb.r == nil {
		rb.r = bytes.NewReader(utils.MakeJsonBodyRaw(rb.Body))
	}
	return rb.r
}

func (rb *RequestBodyJSON[BodyT]) BodyRaw() []byte {
	return utils.MakeJsonBodyRaw(rb.Body)
}

func (rb *RequestBodyJSON[JsonBodyT]) SetBodyStream(r io.Reader) error {
	// var buf bytes.Buffer
	// io.TeeReader(r, &buf)
	data, err := io.ReadAll(r)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, &rb.Body)
}

type RequestBodyStream struct {
	r io.Reader
}

func (rb *RequestBodyStream) BodyStream() io.Reader {
	if rb.r == nil {
		// for safety and convenience
		// (otherwise, calling side always has to check that BodyStream() != nil)
		rb.r = bytes.NewReader(nil)
	}
	return rb.r
}

func (rb *RequestBodyStream) SetBodyStream(r io.Reader) error {
	rb.r = r
	return nil
}

// bound

type requestMethodBound struct{}

func (r requestMethodBound) SetMethod(string) error {
	return errBoundRequestMethod
}

var errBoundRequestMethod = errors.New("request method is bound to request type and can't be set")

// DEPRECATED
// TODO: remove
type RequestBase[UriT any, QueryT any, BodyT any] struct {
	URI   UriT
	Query QueryT
	Body  BodyT
}
