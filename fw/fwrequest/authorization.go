package fwrequest

// TODO: add header preset for authorization headers

// HavingAuthToken sets token to proper place (e.g. Authorization header).
type HavingAuthToken interface {
	AuthToken() string
	SetAuthToken(string) error
}
