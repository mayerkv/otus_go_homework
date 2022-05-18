package main

import (
	"io/fs"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCopy(t *testing.T) {
	t.Run("when source file not exists, then returns error", fileNotExists)

	t.Run("when source is directory, then returns error", sourceIsDirectory)

	t.Run("when offset grate then file size", offsetGrateThenFileSize)
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
