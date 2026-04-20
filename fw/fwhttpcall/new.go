package fwhttpcall

import (
	"context"

	"github.com/Deimvis-go/fw/fw/fwhttp"
	"github.com/Deimvis-go/fw/fw/internal/types"
	"github.com/Deimvis/go-ext/go1.25/xchaincall/xwrapcall"
)

// TODO: impl
// Want something like: fwhttpcall.
// 	New(rs, req).
//  With(
//    fwhttpcallmw.ResponseValidation()
//    fwhttpcallmw.RequestValidation()
//    fwhttpcallmw.RequestGlobalEndpointFromURL(baseUrl),
//    fwhttpcallmw.MoveResponseToFirstMatched(fwmatch.ByCode, &respCandidate{})
//  ).
//  Do(???)
//  Do(fwhttpcall.ActionFromRequester(fwhttp.Requester)) ???
//  Do(fwhttpcall.RequesterCallAction):
//  	Do(func(c Context) error { resp, ... := c.Requester.Do(c.Request()) ; c.SetResponse(resp) ; c.SetBodyCloseFn(close) ; ... }
//
// upd: I consider fwhttpcall should hide xwrapcall, so correct builder will be:
//    resp, close, err := New(rs, req).With(...).Do(...opts)

type Context interface {
	context.Context
	// ???
	Requester() fwhttp.Requester
	Request() types.Request
	Response() types.Response
	SetResponse(types.Response)
	SetBodyCloseFn(fwhttp.BodyCloseFn)
}

var _ xwrapcall.Context = (Context)(nil)
