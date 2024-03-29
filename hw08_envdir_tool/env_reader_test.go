package main

import (
	"io/fs"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadDir(t *testing.T) {
	pathErr := func(t *testing.T, err error) {
		t.Helper()
		var exp *fs.PathError
		require.ErrorAs(t, err, &exp)
	}

	expEnvironments := Environment{
		"BAR":   NewEnv("bar"),
		"EMPTY": NewEnv(" "),
		"HELLO": NewEnv(`"hello"`),
		"UNSET": NewEnv(""),
		"FOO":   NewEnv("   foo\nwith new line"),
	}

	tests := []struct {
		name       string
		checkError func(t *testing.T, err error)
		dirName    string
		withEnv    bool
	}{
		{
			name:       "dir not exist",
			checkError: pathErr,
			dirName:    "ada",
		},
		{
			name:       "is not dir",
			checkError: pathErr,
			dirName:    "env_reader.go",
		},
		{
			name: "value is invalid",
			checkError: func(t *testing.T, err error) {
				t.Helper()
				require.ErrorIs(t, err, ErrValueInvalid)
			},
			dirName: "testdata",
		},
		{
			name:    "success",
			dirName: "testdata/env",
			withEnv: true,
			checkError: func(t *testing.T, err error) {
				t.Helper()
				require.Nil(t, err)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ReadDir(tt.dirName)
			tt.checkError(t, err)
			if !tt.withEnv {
				require.Nil(t, got)
			} else {
				require.Equal(t, len(expEnvironments), len(got))
				for k, expVal := range expEnvironments {
					gotVal, ok := got[k]
					require.True(t, ok, k)
					require.Equal(t, expVal, gotVal, k)
				}
			}
		})
	}
}
