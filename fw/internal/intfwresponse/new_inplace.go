package intfwresponse

import (
	"io"
	"net/http"

	"github.com/Deimvis-go/fw/fw/internal/types"
	"github.com/Deimvis/go-ext/go1.25/xhttp"
)

func NewInplace() types.ResponseReflect {
	return &impl{}
}

type impl struct {
	code   int
	header http.Header
	body   io.Reader
}

func (i *impl) Code() int {
	return i.code
}

func (i *impl) SetCode(c int) error {
	i.code = c
	return nil
}

func (i *impl) Header() xhttp.ConstHeader {
	return xhttp.AsConstHeader(i.header)
}

func (i *impl) SetHeader(h http.Header) error {
	i.header = h
	return nil
}

func (i *impl) BodyStream() io.Reader {
	return i.body
}

func (i *impl) SetBodyStream(r io.Reader) error {
	i.body = r
	return nil
}
