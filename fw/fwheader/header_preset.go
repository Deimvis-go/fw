package fwheader

import "net/http"

type Preset interface {
	New() http.Header
	// TODO: maybe add Validate(http.Header) error method for SetHeader validation
}

type NoPreset struct{}

func (p NoPreset) New() http.Header {
	return make(http.Header)
}

func CustomHeaderPreset(h http.Header) Preset {
	return customPreset{h: h}
}

type customPreset struct {
	h http.Header
}

func (p customPreset) New() http.Header {
	return p.h.Clone()
}

// TODO: somehow support Accept header (=application/json)
type JSONPreset struct{}

func (p JSONPreset) New() http.Header {
	// TODO: optimize with pool of preallocated http.Header objects
	h := make(http.Header)
	h["Content-Type"] = jsonContentType
	return h
}
