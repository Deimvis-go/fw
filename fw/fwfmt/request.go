package fwfmt

import (
	"fmt"

	"github.com/Deimvis-go/fw/fw/internal/types"
	"github.com/Deimvis/go-ext/go1.25/xfmt"
)

func Request(r types.Request) string {
	kvs := []any{
		"method", r.Method(),
		"path", r.Path(),
		"query", r.QueryString(),
	}
	if rr, ok := r.(types.BufferedRequest); ok {
		kvs = append(kvs, "body", string(rr.BodyRaw()))
	}
	return fmt.Sprintf("request(%s)", xfmt.Sprintfkv("%v=%v", ", ", kvs...))
}
