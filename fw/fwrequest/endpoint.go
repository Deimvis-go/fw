package fwrequest

import (
	"github.com/Deimvis-go/fw/fw/internal/types"
	"github.com/Deimvis/go-ext/go1.25/xnet/xurl"
)

type LocalEndpoint interface {
	Method() string
	// PathTemplate also serves a role of canonical path.
	PathTemplate(opts ...PathTemplateGenerationOption) PathTemplate

	// TODO: request element like this:
	// fw.RequestPath[fwpath.ConstTemplate("/user/{id}", path.BraceEscaping)]
	// ^^^ how to push template value into type definition?
	// fw.RequestPath[fwpath.Template{Const: "/user/{id}", Format: path.BraceEscaping}]
	// this element will implement methods PathTemplate()
	// and this element will give method pathFormat() and so to allow easily implement Path() string method.
	// Even better, Path() string could be implemented under the hood if we define uri params order somehow in tags.
}

type GlobalEndpoint interface {
	LocalEndpoint
	Scheme() string
	Authority() xurl.Authority
}

// NewEndpointURITemplate constructs template for uri
// which serves a role of a canonical uri
// and can be used to check endpoints equality.
// If input implements GlobalEndpoint,
// then GlobalEndpoint info is used.
func NewEndpointURITemplate(le LocalEndpoint) {
	// Cast to GlobalEndpoint if possible
	panic("TODO: implement")
}

// interface shortcuts

type RequestWithLocalEndpoint interface {
	types.Request
	LocalEndpoint
}

type RequestWithGlobalEndpoint interface {
	types.Request
	GlobalEndpoint
}
