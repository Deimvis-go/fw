package fwss

import (
	"github.com/Deimvis-go/fw/fw"
	"github.com/Deimvis-go/fw/fw/fwresponse"
	"github.com/Deimvis/go-ext/go1.25/xfmt"
)

func Resp400(msgAndArgs ...any) *fw.ErrorResponse400NoHeader {
	resp := &fw.ErrorResponse400NoHeader{}
	resp.Body.Error = xfmt.Sprintfg(msgAndArgs...)
	return resp
}

func Resp401(msgAndArgs ...any) *fw.ErrorResponse401NoHeader {
	resp := &fw.ErrorResponse401NoHeader{}
	resp.Body.Error = xfmt.Sprintfg(msgAndArgs...)
	return resp
}

func Resp403(msgAndArgs ...any) *fw.ErrorResponse403NoHeader {
	resp := &fw.ErrorResponse403NoHeader{}
	resp.Body.Error = xfmt.Sprintfg(msgAndArgs...)
	return resp
}

func Resp404(msgAndArgs ...any) *fw.ErrorResponse404NoHeader {
	resp := &fw.ErrorResponse404NoHeader{}
	resp.Body.Error = xfmt.Sprintfg(msgAndArgs...)
	return resp
}

func Resp409(msgAndArgs ...any) *fw.ErrorResponse409NoHeader {
	resp := &fw.ErrorResponse409NoHeader{}
	resp.Body.Error = xfmt.Sprintfg(msgAndArgs...)
	return resp
}

func Resp500(msgAndArgs ...any) *fw.ErrorResponse500NoHeader {
	resp := &fw.ErrorResponse500NoHeader{}
	resp.Body.Error = xfmt.Sprintfg(msgAndArgs...)
	return resp
}

func Resp503(msgAndArgs ...any) *fw.ErrorResponse503NoHeader {
	resp := &fw.ErrorResponse503NoHeader{}
	resp.Body.Error = xfmt.Sprintfg(msgAndArgs...)
	return resp
}

func ErrorResp(err error) *fwresponse.ErrorResponse {
	resp := &fwresponse.ErrorResponse{}
	resp.Body.Error = err.Error()
	return resp
}
