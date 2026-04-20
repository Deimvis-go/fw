package fwheader

import "net/http"

type HavingInternals interface {
	// HeaderDirect returns http.Header unlike Header() returning xhttp.ConstHeader.
	// It may return a mutable copy of request/response header,
	// so changes to the returned http.Header MAY affect request/response state.
	// HeaderDirect is used for efficient moving of underlying header,
	// liminating extra copying imposed by using xhttp.ConstHeader.
	HeaderDirect() http.Header
}
