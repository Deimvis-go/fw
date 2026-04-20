package fw

import (
	"bytes"
	"io"

	"github.com/Deimvis/go-ext/go1.25/xcheck/xmust"
)

func MustBufferizeRequest(r Request) BufferedRequest {
	return xmust.Do(BufferizeRequest(r))
}

func MustBufferizeResponse(r Response) BufferedResponse {
	return xmust.Do(BufferizeResponse(r))
}

func BufferizeRequest(r Request) (BufferedRequest, error) {
	if rb, ok := r.(BufferedRequest); ok {
		// already buffered
		return rb, nil
	}

	body, err := io.ReadAll(r.BodyStream())
	if err != nil {
		return nil, err
	}
	return bufferedRequest{Request: r, buf: body}, nil
}

func BufferizeResponse(r Response) (BufferedResponse, error) {
	if rb, ok := r.(BufferedResponse); ok {
		// already buffered
		return rb, nil
	}

	body, err := io.ReadAll(r.BodyStream())
	if err != nil {
		return nil, err
	}
	return bufferedResponse{Response: r, buf: body}, nil
}

type bufferedRequest struct {
	Request
	buf []byte
}

func (r bufferedRequest) BodyRaw() []byte {
	return r.buf
}

// NOTE: this method is overriden because original BodyStream()
// was completely read on bufferedRequest creation.
func (r bufferedRequest) BodyStream() io.Reader {
	return bytes.NewReader(r.buf)
}

type bufferedResponse struct {
	Response
	buf []byte
}

func (r bufferedResponse) BodyRaw() []byte {
	return r.buf
}

// NOTE: this method is overriden because original BodyStream()
// was completely read on bufferedResponse creation.
func (r bufferedResponse) BodyStream() io.Reader {
	return bytes.NewReader(r.buf)
}
