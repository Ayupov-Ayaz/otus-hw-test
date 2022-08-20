package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

const (
	inFileName  = "./testdata/input.txt"
	outFileName = "output.txt"
)

func TestCopy(t *testing.T) {
	checkSize := func(size int64) func(t *testing.T) {
		return func(t *testing.T) {
			info, err := os.Stat(outFileName)
			require.Nil(t, err)
			require.Equal(t, size, info.Size())
		}
	}

	fileNotFound := func(t *testing.T) {
		_, err := os.Stat(outFileName)
		require.ErrorIs(t, err, os.ErrNotExist)
	}
	tests := []struct {
		name   string
		limit  int64
		offset int64
		after  func(t *testing.T)
		err    error
	}{
		{
			name:   "offset > size",
			offset: 6618,
			err:    ErrOffsetExceedsFileSize,
			after:  fileNotFound,
			limit:  1,
		},
		{
			name:  "copy 0 byte",
			after: fileNotFound,
		},
		{
			name:  "copy 10 byte",
			limit: 10,
			after: checkSize(10),
		},
		{
			name:  "copy full file",
			limit: 7000,
			after: checkSize(6617),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Copy(inFileName, outFileName, tt.offset, tt.limit)
			require.ErrorIs(t, err, tt.err)
			err = os.Remove(outFileName)
			require.Nil(t, err)
		})
	}
}

func Test_needRead(t *testing.T) {
	const size int64 = 6617

	tests := []struct {
		offset int64
		limit  int64
		exp    int64
	}{
		{
			offset: 0,
			limit:  0,
			exp:    size,
		},
		{
			offset: 0,
			limit:  10,
			exp:    10,
		},
		{
			offset: 0,
			limit:  1000,
			exp:    1000,
		},
		{
			offset: 0,
			limit:  10000,
			exp:    size,
		},
		{
			offset: 100,
			limit:  1000,
			exp:    1000,
		},
		{
			offset: 6000,
			limit:  1000,
			exp:    size - 6000,
		},
	}

	for _, tt := range tests {
		got := needRead(size, tt.offset, tt.limit)
		require.Equal(t, tt.exp, got)
	}
}
