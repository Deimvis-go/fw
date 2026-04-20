package fw

import (
	"bytes"
	"io"
	"net/http"

	"github.com/Deimvis-go/fw/fw/fwheader"
	"github.com/Deimvis-go/fw/fw/fwresponse"
	"github.com/Deimvis-go/fw/fw/internal/intfwresponse"
	"github.com/Deimvis/go-ext/go1.25/xhttp"
)

// TODO: move elements to separate package for clarity

// Code

type Response200 struct {
	responseCodeBound
}

func (r *Response200) Code() int {
	return 200
}

type Response202 struct {
	responseCodeBound
}

func (r *Response202) Code() int {
	return 202
}

type Response204 struct {
	responseCodeBound
}

func (r *Response204) Code() int {
	return 204
}

type Response302 struct {
	responseCodeBound
}

func (r *Response302) Code() int {
	return 302
}

type Response400 struct {
	responseCodeBound
}

func (r *Response400) Code() int {
	return 400
}

type Response401 struct {
	responseCodeBound
}

func (r *Response401) Code() int {
	return 401
}

type Response403 struct {
	responseCodeBound
}

func (r *Response403) Code() int {
	return 403
}

type Response404 struct {
	responseCodeBound
}

func (r *Response404) Code() int {
	return 404
}

type Response409 struct {
	responseCodeBound
}

func (r *Response409) Code() int {
	return 409
}

type Response500 struct {
	responseCodeBound
}

func (r *Response500) Code() int {
	return 500
}

type Response503 struct {
	responseCodeBound
}

func (r Response503) Code() int {
	return 503
}

// Header

type ResponseNoHeader struct{}

func (r *ResponseNoHeader) Header() xhttp.ConstHeader {
	// not returning nil, but const header wrap,
	// because it's safe to call methods over it.
	return xhttp.AsConstHeader(nil)
}

func (r *ResponseNoHeader) SetHeader(http.Header) error {
	// ignore extra fields
	return nil
}

// TODO: migrate to ResponseStructHeader and rename ResponseStructHeader to ResponseHeader
type ResponseHeader[PresetT HeaderPreset] struct {
	// TODO: maybe make Headers unexported
	Headers http.Header `json:"-"`
	hp      PresetT
}

func (r *ResponseHeader[PresetT]) Header() xhttp.ConstHeader {
	if r.Headers == nil {
		r.Headers = r.hp.New()
	}
	return xhttp.AsConstHeader(r.Headers)
}

func (r *ResponseHeader[PresetT]) SetHeader(h http.Header) error {
	r.Headers = h
	return nil
}

// GOEXPERIMENT=aliastypeparams : make alias
// type ResponseStructHeader[PresetT HeaderPreset, T any] = fwheader.Structured[PresetT, T]
type ResponseStructHeader[PresetT HeaderPreset, T any] struct {
	fwheader.Structured[PresetT, T]
}

// Body

type ResponseNoBody struct{}

func (rb *ResponseNoBody) BodyStream() io.Reader {
	return bytes.NewReader(nil)
}

func (rb *ResponseNoBody) BodyRaw() []byte {
	return nil
}

func (rb *ResponseNoBody) SetBodyStream(io.Reader) error {
	// ignore extra fields
	return nil
}

func (rb *ResponseNoBody) SetBodyRaw(content []byte) error {
	// ignore extra fields
	return nil
}

// GOEXPERIMENT=aliastypeparams : make alias
// type ResponseBodyJSON[BodyT any] = fwresponse.BodyJSON[BodyT]
type ResponseBodyJSON[BodyT any] struct {
	fwresponse.BodyJSON[BodyT]
}

type ResponseBodyStream struct {
	r io.Reader
}

func (rb *ResponseBodyStream) BodyStream() io.Reader {
	return rb.r
}

func (rb *ResponseBodyStream) SetBodyStream(r io.Reader) error {
	rb.r = r
	return nil
}

// bound

type responseCodeBound struct{}

func (r responseCodeBound) SetCode(int) error {
	return intfwresponse.ErrBoundResponseCode
}
