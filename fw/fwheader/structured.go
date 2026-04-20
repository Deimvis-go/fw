package fwheader

import (
	"fmt"
	"net/http"

	"github.com/Deimvis-go/fw/fw/internal/headerenc"
	"github.com/Deimvis-go/fw/fw/internal/utils"
	"github.com/Deimvis-go/valid"
	"github.com/Deimvis/go-ext/go1.25/xhttp"
)

type Structured[PresetT Preset, T any] struct {
	Headers T
	hp      PresetT
}

// Header returns a copy of current state as xhttp.ConstHeader.
func (s *Structured[PresetT, T]) Header() xhttp.ConstHeader {
	h := s.hp.New()
	sh := utils.MakeHttpHeader(s.Headers)
	for k, vs := range sh {
		for _, v := range vs {
			h.Add(k, v)
		}
	}
	if havingExtras, ok := any(s.Headers).(HavingExtras); ok {
		extras := havingExtras.Extra()
		for k, vs := range extras {
			if _, ok := sh[k]; ok {
				panic(fmt.Errorf("extra headers overlap with structured headers on key `%s`", k))
			}
			for _, v := range vs {
				h.Add(k, v)
			}
		}
	}
	return xhttp.AsConstHeader(h)
}

func (s *Structured[PresetT, T]) SetHeader(h http.Header) error {
	var extras http.Header = nil
	if havingExtras, ok := any(s.Headers).(HavingExtras); ok {
		extras = havingExtras.Extra()
	}
	err := headerenc.Unmarshal(xhttp.AsConstHeader(h), &s.Headers, extras)
	if err != nil {
		return err
	}
	return valid.Deep(s.Headers)
}
