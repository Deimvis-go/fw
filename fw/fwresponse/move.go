package fwresponse

import (
	"github.com/Deimvis-go/fw/fw/fwheader"
	"github.com/Deimvis-go/fw/fw/internal/types"
)

func Move(src types.Response, dst types.ResponseReflect) error {
	dst.SetCode(src.Code())
	if srcInt, ok := src.(fwheader.HavingInternals); ok {
		dst.SetHeader(srcInt.HeaderDirect())
	} else {
		dst.SetHeader(src.Header().Clone())
	}
	dst.SetBodyStream(src.BodyStream())
	return nil
}
