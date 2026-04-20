package fwhttp

type BodyCloseFn = func() error

var (
	NoopBodyCloseFn BodyCloseFn = func() error { return nil }
)
