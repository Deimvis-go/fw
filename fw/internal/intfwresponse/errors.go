package intfwresponse

import "errors"

var ErrBoundResponseCode = errors.New("response code is bound to request type and can't be set")
