package fwfmt

import (
	"fmt"

	"github.com/Deimvis-go/fw/fw/internal/types"
	"github.com/Deimvis/go-ext/go1.25/xfmt"
)

// TODO: options to customize fields

func Response(r types.Response) string {
	kvs := []any{
		"code", r.Code(),
	}
	if rr, ok := r.(types.BufferedResponse); ok {
		kvs = append(kvs, "body", string(rr.BodyRaw()))
	}
	return fmt.Sprintf("response(%s)", xfmt.Sprintfkv("%v=%v", ", ", kvs...))
}
