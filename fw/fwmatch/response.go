package fwmatch

import (
	"github.com/Deimvis-go/fw/fw/internal/types"
)

// TODO: build matches on top of (xcheck/) xwhether predicates

type ResponsePred func(v1, v2 types.Response) bool

func ByCode(v1, v2 types.Response) bool {
	return v1.Code() == v2.Code()
}

var _ ResponsePred = ByCode
