package utils

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestToString(t *testing.T) {
	type isPanic bool
	var panics, nopanic isPanic = true, false

	tcs := []struct {
		title    string
		v        any
		exp      string
		expPanic isPanic
	}{
		{
			"string",
			"hello world !?<=>(-)*:#",
			"hello world !?<=>(-)*:#",
			nopanic,
		},
		{
			"alised string",
			strAlias,
			"hello world !?<=>(-)*:#",
			nopanic,
		},
		{
			"custom string",
			strCustom,
			"hello world !?<=>(-)*:#",
			nopanic,
		},
		{
			"struct",
			struct{}{},
			"",
			panics,
		},
	}
	for _, tc := range tcs {
		t.Run(tc.title, func(t *testing.T) {
			panicked := false
			func() {
				if tc.expPanic {
					defer func() {
						r := recover()
						panicked = r != nil
					}()
				}
				act := toString(reflect.ValueOf(tc.v))
				require.Equal(t, tc.exp, act)
			}()
			if tc.expPanic {
				require.True(t, panicked)
			}
		})
	}
}

type StrAlias = string
type StrCustom string

var (
	strAlias  StrAlias  = "hello world !?<=>(-)*:#"
	strCustom StrCustom = "hello world !?<=>(-)*:#"
)
