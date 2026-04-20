package fwhttp

import (
	"github.com/Deimvis-go/fw/fw/fwmatch"
	"github.com/Deimvis-go/fw/fw/internal/types"
	"github.com/Deimvis/go-ext/go1.25/xhttp"
)

// TODO: WithMoveToMatched(matcher fwmatch.ResponseReflectMatcher)
// WithMoveToFirstMatched(pred) = WithMoveToMatched(fwmatch.First(pred))

// NOTE: Only WithmoveToFirstMatched is implemented, because we want statically check
// that response candidates implement ResponseReflect.
// It's not necessary for matching, but neccesary for moving.

// WithMoveToFirstMatched matches response and moves content to first matched.
// If no response is matched, NoResponseMatch error is returned.
func WithMoveToFirstMatched(pred fwmatch.ResponsePred, opts ...types.ResponseReflect) DoOption {
	matcher := func(r types.Response) (types.Response, bool) {
		for _, opt := range opts {
			if pred(r, opt) {
				return opt, true
			}
		}
		return nil, false
	}
	return func(c *doCfg) {
		c.matcherForMove = matcher
	}
}

// TODO: it's a hack, get rid of it in favor of Extraheaders
// working only when fw.Request's header implements fwheader.Overridable
func WithExtraHeaders(h xhttp.ConstHeader) DoOption {
	return func(c *doCfg) {
		c.extraHeaders = h
	}
}

type doCfg struct {
	matcherForMove fwmatch.ResponseMatcher
	extraHeaders   xhttp.ConstHeader
}
