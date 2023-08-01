package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// func TestMain(m *testing.M) {
// 	err := checkParams()
// 	require.ErrorIs(m, err, errInvalidArgs, "Errors are not the same. Actual error is %v", err)
// 	os.Exit(m.Run())
// }

func TestValidArgs(t *testing.T) {
	t.Run("invalid", func(t *testing.T) {
		err := checkParams()
		require.ErrorIs(t, err, errInvalidArgs, "Errors are nont the same. Actual error is %v", err)
	})
}
