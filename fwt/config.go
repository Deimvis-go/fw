package fwt

import (
	"net/http"

	"go.uber.org/zap"
)

type Config struct {
	SERVER       *http.Server
	ACCESS_TOKEN *string
	LOGGER       *zap.SugaredLogger
}

func SetConfigRefreshFunc(fn func(c *Config)) {
	refreshFn = fn
}

func WithConfigRefreshFunc(fn func(c *Config)) func() {
	origFn := refreshFn
	refreshFn = fn
	return func() {
		refreshFn = origFn
	}
}

var cfg = &Config{}
var refreshFn func(c *Config)
