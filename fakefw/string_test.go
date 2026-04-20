package fakefw

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestString_Deterministic(t *testing.T) {
	require.Equal(t, "FxhfnEwYfjyic61N", String("request_id"))
	require.Equal(t, "FxhfnEwYfjyic61N", String("request_id"))
}
