package fwmatch

import (
	"github.com/Deimvis-go/fw/fw/internal/types"
)

type ResponseMatcher func(types.Response) (types.Response, bool)

func First[T any, U any, Pred func(T, U) bool](pred Pred, opts ...U) func(T) (U, bool) {
	var emptyMatch U
	fn := func(r T) (U, bool) {
		for _, opt := range opts {
			if pred(r, opt) {
				return opt, true
			}
		}
		return emptyMatch, false
	}
	return fn
}
