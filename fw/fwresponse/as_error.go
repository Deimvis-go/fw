package fwresponse

func AsError(he ResponseHavingError) error {
	return asError{ResponseHavingError: he}
}

type asError struct {
	ResponseHavingError
}

func (ae asError) Error() string {
	return ae.ResponseHavingError.Error()
}
