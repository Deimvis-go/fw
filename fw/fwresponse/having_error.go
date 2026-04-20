package fwresponse

import "github.com/Deimvis-go/fw/fw/internal/types"

type HavingError interface {
	Error() string
}

type HavingError_hidden interface {
	error_() string
}

func RevealResponseHavingError(r ResponseHavingError_hidden) ResponseHavingError {
	return withErrorRevealed{ResponseHavingError_hidden: r}
}

type withErrorRevealed struct {
	ResponseHavingError_hidden
}

func (her withErrorRevealed) Error() string {
	return her.ResponseHavingError_hidden.error_()
}

// interface shortcuts

type ResponseHavingError interface {
	types.Response
	HavingError
}

type ResponseHavingError_hidden interface {
	types.Response
	HavingError_hidden
}
