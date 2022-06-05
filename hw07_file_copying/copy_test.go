package main

import (
	"bytes"
	"io/fs"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCopy(t *testing.T) {
	t.Run("when source file not exists, then returns error", fileNotExists)

	t.Run("when destination path is empty, then returns error", destinationPathIsEmpty)

	t.Run("when source is directory, then returns error", sourceIsDirectory)

	t.Run("when offset grate then file size", offsetGrateThenFileSize)

	t.Run("cases", cases)
}

func destinationPathIsEmpty(t *testing.T) {
	err := Copy("testdata/input.txt", "", 0, 0)

	require.Error(t, err)
	require.IsType(t, &fs.PathError{}, err)
}

func fileNotExists(t *testing.T) {
	err := Copy("non_existent_file", "/dev/null", 0, 0)

	require.Error(t, err)
	require.IsType(t, &fs.PathError{}, err)
}

func sourceIsDirectory(t *testing.T) {
	err := Copy("testdata", "/dev/null", 0, 0)

	require.Error(t, err)
	require.ErrorIs(t, ErrUnsupportedFile, err)
}

func offsetGrateThenFileSize(t *testing.T) {
	var offset int64 = 7000
	err := Copy("testdata/input.txt", "/dev/null", offset, 0)

	require.Error(t, err)
	require.ErrorIs(t, ErrOffsetExceedsFileSize, err)
}

func cases(t *testing.T) {
	tests := []struct {
		name string
		args struct {
			to     string
			limit  int64
			offset int64
		}
		want string
	}{
		{
			name: "offset 0 limit 0",
			args: struct {
				to     string
				limit  int64
				offset int64
			}{to: "/tmp/out_offset0_limit0.txt", limit: 0, offset: 0},
			want: "testdata/out_offset0_limit0.txt",
		},
		{
			name: "offset 0 limit 10",
			args: struct {
				to     string
				limit  int64
				offset int64
			}{to: "/tmp/out_offset0_limit10.txt", limit: 10, offset: 0},
			want: "testdata/out_offset0_limit10.txt",
		},
		{
			name: "offset 0 limit 1000",
			args: struct {
				to     string
				limit  int64
				offset int64
			}{to: "/tmp/out_offset0_limit1000.txt", limit: 1000, offset: 0},
			want: "testdata/out_offset0_limit1000.txt",
		},
		{
			name: "offset 0 limit 10000",
			args: struct {
				to     string
				limit  int64
				offset int64
			}{to: "/tmp/out_offset0_limit10000.txt", limit: 10000, offset: 0},
			want: "testdata/out_offset0_limit10000.txt",
		},
		{
			name: "offset 100 limit 1000",
			args: struct {
				to     string
				limit  int64
				offset int64
			}{to: "/tmp/out_offset100_limit1000.txt", limit: 1000, offset: 100},
			want: "testdata/out_offset100_limit1000.txt",
		},
		{
			name: "offset 6000 limit 1000",
			args: struct {
				to     string
				limit  int64
				offset int64
			}{to: "/tmp/out_offset6000_limit1000.txt", limit: 1000, offset: 6000},
			want: "testdata/out_offset6000_limit1000.txt",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if err := Copy("testdata/input.txt", test.args.to, test.args.offset, test.args.limit); err != nil {
				t.Error(err)
			}

			expected, err := ioutil.ReadFile(test.want)
			if err != nil {
				t.Error(err)
			}
			actual, err := ioutil.ReadFile(test.args.to)
			if err != nil {
				t.Error(err)
			}

			require.True(t, bytes.Equal(expected, actual))
		})
	}
}
