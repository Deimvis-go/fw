package fw

import (
	"context"
	"net/http"

	"github.com/Deimvis-go/fw/fw/fwrequest"
	"github.com/Deimvis-go/fw/fw/fwresponse"
)

// CopyToHttpRequest constructs *http.Request from a copy of BufferedRequest data.
func CopyToHttpRequest(r BufferedRequest) (*http.Request, error) {
	return fwrequest.CopyToHttp(r)
}

// CopyToHttpRequestWithContext constructs *http.Request from a copy of BufferedRequest data.
func CopyToHttpRequestWithContext(ctx context.Context, r BufferedRequest) (*http.Request, error) {
	return fwrequest.CopyToHttpWithContext(ctx, r)
}

// MoveToHttpRequest constructs *http.Request reusing Request data as much as possible.
// Body stream is always reused in returned *http.Request.
func MoveToHttpRequest(r Request) (*http.Request, error) {
	return fwrequest.MoveToHttp(r)
}

// MoveToHttpRequestWithContext constructs *http.Request reusing Request data as much as possible.
// Body stream is always reused in returned *http.Request.
func MoveToHttpRequestWithContext(ctx context.Context, r Request) (*http.Request, error) {
	return fwrequest.MoveToHttpWithContext(ctx, r)
}

// MoveFromHttpResponse constructs fw.Response reusing *http.Request data as much as possible.
// Note that *http.Response body closing is on the calling side.
func MoveFromHttpResponse(httpResp *http.Response, r ResponseReflect) error {
	return fwresponse.MoveFromHttp(httpResp, r)
}
