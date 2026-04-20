package fwhttp

import (
	"context"

	"github.com/Deimvis-go/fw/fw/internal/types"
)

type Requester interface {
	Do(context.Context, types.Request, ...DoOption) (types.Response, BodyCloseFn, error)
}
