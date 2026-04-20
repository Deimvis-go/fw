package fwresponse

type ErrorResponse struct {
	BodyJSON[struct {
		Error string `json:"error"`
	}]
}

var _ HavingError_hidden = ErrorResponse{}

func (er ErrorResponse) error_() string {
	return er.Body.Error
}
