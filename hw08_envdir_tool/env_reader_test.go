package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadDir(t *testing.T) {
	t.Run("Read dir eror", func(t *testing.T) {
		env, err := ReadDir("C:\\ProgramFiles\\some dir")
		require.Nil(t, env, "actual env is %v", env)
		require.NotEqual(t, nil, err, "actual error is %v", err)
	})

	t.Run("Open file error", func(t *testing.T) {
		f, err := os.Create("testdata/env/perm.txt")
		defer os.Remove(f.Name())
		if err != nil {
			return
		}
		os.Chmod(f.Name(), 0o222)
		env, err := ReadDir("testdata/env")
		require.Nil(t, env, "actual env is %v", env)
		require.NotEqual(t, nil, err, "actual error is %v", err)
	})

	t.Run("Success!", func(t *testing.T) {
		expectedEnv := Environment{
			"BAR":   EnvValue{"bar", false},
			"EMPTY": EnvValue{"", false},
			"FOO":   EnvValue{"   foo\nwith new line", false},
			"HELLO": EnvValue{"\"hello\"", false},
			"UNSET": EnvValue{"", true},
		}
		env, err := ReadDir("testdata/env")
		require.Nil(t, err, "actual error is %v", err)
		require.Equal(t, expectedEnv, env, "actual env is %v", env)
	})
}
