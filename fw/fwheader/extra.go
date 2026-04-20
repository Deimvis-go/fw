package fwheader

import "net/http"

type HavingExtras interface {
	Extra() http.Header
}

type WithExtras struct {
	extra *http.Header
}

var _ HavingExtras = &WithExtras{}

func (e WithExtras) Extra() http.Header {
	if e.extra == nil {
		h := make(http.Header)
		e.extra = &h
	}
	return *e.extra
}
