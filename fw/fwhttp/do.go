package fwhttp

import (
	"context"

	"github.com/Deimvis-go/fw/fw/fwmatch"
	"github.com/Deimvis-go/fw/fw/fwrequest"
	"github.com/Deimvis-go/fw/fw/fwresponse"
	"github.com/Deimvis-go/fw/fw/internal/intfwresponse"
	"github.com/Deimvis-go/fw/fw/internal/types"
	"github.com/Deimvis/go-ext/go1.25/ext"
	"github.com/Deimvis/go-ext/go1.25/xcheck/xmust"
	"github.com/Deimvis/go-ext/go1.25/xhttp"
)

type DoOption = func(*doCfg)

func Do(ctx context.Context, rs xhttp.Requester, req fwrequest.RequestWithGlobalRoute, opts ...DoOption) (types.Response, BodyCloseFn, error) {
	var cfg doCfg
	for _, opt := range opts {
		opt(&cfg)
	}

	// note: always return non-nil body close fn for safety
	httpReq, err := fwrequest.MoveToHttpWithContext(ctx, req)
	if err != nil {
		return nil, nil, err
	}
	if cfg.extraHeaders != nil {
		cfg.extraHeaders.Range(func(k string, vs []string) bool {
			for _, v := range vs {
				httpReq.Header.Add(k, v)
			}
			return true
		})
	}
	httpResp, err := rs.Do(httpReq)
	if err != nil {
		return nil, NoopBodyCloseFn, err
	}
	close := httpResp.Body.Close
	defer ext.OnPanic(func(_ any) {
		_ = close() // ignore error
	})
	resp := intfwresponse.NewInplace()
	err = fwresponse.MoveFromHttp(httpResp, resp)
	if err != nil {
		close()
		return nil, NoopBodyCloseFn, err
	}

	if cfg.matcherForMove != nil {
		matchedResp, ok := cfg.matcherForMove(resp)
		if !ok {
			return nil, NoopBodyCloseFn, fwmatch.NewNoResponseMatch(resp)
		}
		matchedRespMutable, ok := matchedResp.(types.ResponseReflect)
		xmust.True(ok, "matched response do not implement ResponseReflect")
		fwresponse.Move(resp, matchedRespMutable)
		resp = matchedRespMutable
	}

	return resp, close, nil
}
