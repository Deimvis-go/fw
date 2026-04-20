package fwresponse

import (
	"bytes"
	"encoding/json"
	"io"

	"github.com/Deimvis-go/fw/fw/internal/utils"
)

type BodyJSON[BodyT any] struct {
	Body BodyT
	r    io.Reader
}

func (rb *BodyJSON[BodyT]) BodyStream() io.Reader {
	if rb.r == nil {
		rb.r = bytes.NewReader(utils.MakeJsonBodyRaw(rb.Body))
	}
	return rb.r
}

func (rb *BodyJSON[BodyT]) BodyRaw() []byte {
	return utils.MakeJsonBodyRaw(rb.Body)
}

func (rb *BodyJSON[BodyT]) SetBodyStream(r io.Reader) error {
	// var buf bytes.Buffer
	// io.TeeReader(r, &buf)
	data, err := io.ReadAll(r)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, &rb.Body)
}
