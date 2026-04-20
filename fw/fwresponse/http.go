package fwresponse

import (
	"fmt"
	"net/http"

	"github.com/Deimvis-go/fw/fw/internal/intfwresponse"
	"github.com/Deimvis-go/fw/fw/internal/types"
)

func MoveFromHttp(httpResp *http.Response, r types.ResponseReflect) error {
	var err error
	err = r.SetCode(httpResp.StatusCode)
	if err != nil {
		if err != intfwresponse.ErrBoundResponseCode {
			return err
		}
		if r.Code() != httpResp.StatusCode {
			return fmt.Errorf("bound code %d, but got %d", r.Code(), httpResp.StatusCode)
		}
	}
	err = r.SetHeader(httpResp.Header)
	if err != nil {
		return err
	}
	err = r.SetBodyStream(httpResp.Body)
	if err != nil {
		return err
	}
	return nil
}
