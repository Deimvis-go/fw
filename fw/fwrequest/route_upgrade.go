package fwrequest

import (
	"github.com/Deimvis-go/fw/fw/internal/types"
	"github.com/Deimvis/go-ext/go1.25/xnet/xurl"
)

func UpgradeToGlobalRoute(req types.Request, scheme string, authority xurl.Authority) RequestWithGlobalRoute {
	return globalRouteFromLocal{
		Request:   req,
		scheme:    scheme,
		authority: authority,
	}
}

type globalRouteFromLocal struct {
	types.Request
	scheme    string
	authority xurl.Authority
}

func (gr globalRouteFromLocal) Scheme() string {
	return gr.scheme
}

func (gr globalRouteFromLocal) Authority() xurl.Authority {
	return gr.authority
}
