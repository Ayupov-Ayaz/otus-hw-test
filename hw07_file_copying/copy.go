package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/cheggaaa/pb/v3"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

func closeFile(closer io.Closer) {
	if err := closer.Close(); err != nil {
		log.Println("failed to close file: %w", err)
	}
}

const bufferSize int64 = 1024

func copyTo(reader io.Reader, writer io.Writer, need int64) error {
	bar := pb.StartNew(int(need))
	var done bool
	for !done {
		size := bufferSize
		if need < bufferSize {
			size = need
		}

		data := make([]byte, size)
		_, err := reader.Read(data)
		if err != nil {
			if !errors.Is(err, io.EOF) {
				return err
			}
			done = true
		}

		wrote, err := writer.Write(data)
		if err != nil {
			return err
		}

		need -= int64(wrote)
		bar.Add(wrote)

		if need == 0 {
			done = true
		}
	}

	bar.Finish()

	return nil
}

func checkSizeAndOffset(size, offset int64) error {
	if offset > size {
		return fmt.Errorf("offset = '%d', file size = '%d': %w",
			offset, size, ErrOffsetExceedsFileSize)
	}

	return nil
}

func needRead(size, offset, limit int64) int64 {
	if limit > 0 && limit < size && offset == 0 {
		return limit
	}

	size -= offset
	if limit < size && limit != 0 {
		return limit
	}

	return size
}

func checkFile(file os.FileInfo) error {
	if file.Mode().IsDir() {
		return fmt.Errorf("'%s' is dir: %w",
			file.Name(), ErrUnsupportedFile)
	}

	return nil
}

func Copy(fromPath, toPath string, offset, limit int64) error {
	in, err := os.Open(fromPath)
	if err != nil {
		return err
	}

	defer closeFile(in)

	info, err := in.Stat()
	if err != nil {
		return err
	}

	if err := checkFile(info); err != nil {
		return err
	}

	size := info.Size()

	if err := checkSizeAndOffset(size, offset); err != nil {
		return err
	}

	out, err := os.Create(toPath)
	if err != nil {
		return err
	}

	defer closeFile(out)

	if offset > 0 {
		_, err = in.Seek(offset, io.SeekStart)
		if err != nil {
			return err
		}
	}

	needToRead := needRead(size, offset, limit)

	return copyTo(in, out, needToRead)
}
