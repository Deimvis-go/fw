package fw

import (
	"context"

	"github.com/Deimvis-go/fw/fw/fwhttp"
	"github.com/Deimvis/go-ext/go1.25/xcheck/xmust"
	"github.com/Deimvis/go-ext/go1.25/xhttp"
)

// Do performs http request possibly modifying input request.
// If returned error is not nil, then body is automatically closed
// and no close function returned.
// It uses MoveToHttpRequest in order to cast Request to *http.Request.
// If you need to keep Request unchanged, clone it before call:
//
//	var req fw.Request
//	DoRequest(rs, fw.CloneRequest(req))
//	// req is unchanged
func Do(ctx context.Context, rs xhttp.Requester, req RequestWithGlobalRoute, opts ...DoOption) (Response, fwhttp.BodyCloseFn, error) {
	return fwhttp.Do(ctx, rs, req, opts...)
}

// MustDo is alias to Do that panics if non-nil error returned.
func MustDo(ctx context.Context, rs xhttp.Requester, req RequestWithGlobalRoute, opts ...DoOption) (Response, fwhttp.BodyCloseFn) {
	return xmust.Do2(fwhttp.Do(ctx, rs, req, opts...))
}

type DoOption = fwhttp.DoOption
