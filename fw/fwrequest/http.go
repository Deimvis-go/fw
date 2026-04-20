package fwrequest

import (
	"bytes"
	"context"
	"net/http"

	"github.com/Deimvis-go/fw/fw/fwheader"
	"github.com/Deimvis-go/fw/fw/internal/types"
)

func CopyToHttp(r types.BufferedRequest) (*http.Request, error) {
	res, err := http.NewRequest(r.Method(), NewURL(r).String(), bytes.NewReader(r.BodyRaw()))
	if err != nil {
		return nil, err
	}
	res.Header = r.Header().Clone()
	return res, nil
}

func CopyToHttpWithContext(ctx context.Context, r types.BufferedRequest) (*http.Request, error) {
	res, err := http.NewRequestWithContext(ctx, r.Method(), NewURL(r).String(), bytes.NewReader(r.BodyRaw()))
	if err != nil {
		return nil, err
	}
	res.Header = r.Header().Clone()
	return res, nil
}

func MoveToHttp(r types.Request) (*http.Request, error) {
	res, err := http.NewRequest(r.Method(), NewURL(r).String(), r.BodyStream())
	if err != nil {
		return nil, err
	}
	if headerInts, ok := r.(fwheader.HavingInternals); ok {
		res.Header = headerInts.HeaderDirect()
	} else {
		res.Header = r.Header().Clone()
	}
	return res, nil
}

func MoveToHttpWithContext(ctx context.Context, r types.Request) (*http.Request, error) {
	res, err := http.NewRequestWithContext(ctx, r.Method(), NewURL(r).String(), r.BodyStream())
	if err != nil {
		return nil, err
	}
	if headerInts, ok := r.(fwheader.HavingInternals); ok {
		res.Header = headerInts.HeaderDirect()
	} else {
		res.Header = r.Header().Clone()
	}
	return res, nil
}
