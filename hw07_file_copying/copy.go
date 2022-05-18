package main

import (
	"errors"
	"io"
	"os"

	"github.com/cheggaaa/pb/v3"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	src, err := os.Open(fromPath)
	if err != nil {
		return err
	}
	defer src.Close()

	stat, err := src.Stat()
	if err != nil {
		return err
	}
	if !stat.Mode().IsRegular() {
		return ErrUnsupportedFile
	}
	if offset > stat.Size() {
		return ErrOffsetExceedsFileSize
	}

	maxLimit := stat.Size() - offset
	if limit > maxLimit || limit == 0 {
		limit = maxLimit
	}

	dst, err := os.Create(toPath)
	if err != nil {
		return err
	}
	defer dst.Close()

	if _, err := src.Seek(offset, 0); err != nil {
		return err
	}

	var total, chunkSize int64 = 0, 1 << 20
	if chunkSize > limit {
		chunkSize = limit
	}

	bar := pb.New64(limit)
	bar.Set(pb.Bytes, true)
	bar.Start()
	defer bar.Finish()

	for total < limit {
		written, err := io.CopyN(dst, src, chunkSize)
		bar.Add64(written)
		total += written

		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return err
		}
	}

	return nil
}
