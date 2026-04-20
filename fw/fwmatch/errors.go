package fwmatch

import (
	"fmt"

	"github.com/Deimvis-go/fw/fw/fwfmt"
	"github.com/Deimvis-go/fw/fw/internal/types"
)

func NewNoResponseMatch(r types.Response) NoResponseMatch {
	return NoResponseMatch{r: r}
}

type NoResponseMatch struct {
	r types.Response
}

func (e NoResponseMatch) Response() types.Response {
	return e.r
}

func (e NoResponseMatch) Error() string {
	return fmt.Sprintf("response match not found for %s", fwfmt.Response(e.r))
}
