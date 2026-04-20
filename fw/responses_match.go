package fw

import (
	"fmt"
	"io"
	"net/http"

	"github.com/Deimvis/go-ext/go1.25/xhttp"
)

// TODO: remove this file

// DEPRECATED, use fw.Do(ctx, rs, req, fwhttp.WithMoveToFirstMatch(fwmatch.ByCode))
func MatchResponse(resp *http.Response, opts ...ResponseReflect) (Response, error) {
	for _, opt := range opts {
		if resp.StatusCode == opt.Code() {
			err := opt.SetHeader(resp.Header)
			if err != nil {
				return nil, err
			}
			err = opt.SetBodyStream(resp.Body)
			if err != nil {
				return nil, err
			}
			return opt, nil
		}
	}
	return nil, NewNoResponseMatch(newResponse(resp.StatusCode, resp.Header, resp.Body))
}

func NewNoResponseMatch(resp Response) NoResponseMatch {
	return NoResponseMatch{resp: resp}
}

type NoResponseMatch struct {
	resp Response
}

func (e NoResponseMatch) Response() Response {
	return e.resp
}

func (e NoResponseMatch) Error() string {
	return fmt.Sprintf("response match not found for %s", FormatResponse(e.resp))
}

func newResponse(code int, header http.Header, bodyStream io.Reader) Response {
	return &respImpl{code: code, header: header.Clone(), bodyStream: bodyStream}
}

type respImpl struct {
	code       int
	header     http.Header
	bodyStream io.Reader
}

func (r *respImpl) Code() int {
	return r.code
}

func (r *respImpl) Header() xhttp.ConstHeader {
	return xhttp.AsConstHeader(r.header)
}

func (r *respImpl) BodyStream() io.Reader {
	return r.bodyStream
}
