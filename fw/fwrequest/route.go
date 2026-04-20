package fwrequest

import (
	"net/url"

	"github.com/Deimvis-go/fw/fw/internal/types"
	"github.com/Deimvis/go-ext/go1.25/xnet/xurl"
)

type LocalRoute interface {
	Method() string
	Path() string
	QueryString() string
}

type GlobalRoute interface {
	LocalRoute
	Scheme() string
	Authority() xurl.Authority
}

// NewURL constructs URL out of LocalRoute.
// If input implements GlobalRoute,
// then GlobalRoute info is used.
// Note that fw.Request implements LocalRoute.
func NewURL(lr LocalRoute) *url.URL {
	u := &url.URL{}
	if gr, ok := lr.(GlobalRoute); ok {
		u.Scheme = gr.Scheme()
		u.User = gr.Authority().Userinfo
		u.Host = gr.Authority().Hostport.String()
	}
	u.Path = lr.Path()
	u.RawQuery = lr.QueryString()
	return u
}

// interface shortcuts

type RequestWithLocalRoute interface {
	types.Request
	LocalRoute
}

type RequestWithGlobalRoute interface {
	types.Request
	GlobalRoute
}
