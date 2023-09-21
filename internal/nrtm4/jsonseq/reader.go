package jsonseq

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"os"
	"strings"
)

var RS byte = 0x1E

var ErrEmptyPayload = errors.New("empty payload")
var ErrNotJsonSeq = errors.New("not a JSON seq file")
var ErrExtraneousBytes = errors.New("bytes found before record marker")

type RecordReaderFunc func([]byte, error) error

func ReadStringRecords(jsonSeq string, fn RecordReaderFunc) error {
	reader := bufio.NewReader(strings.NewReader(jsonSeq))
	return ReadRecords(reader, fn)
}

func ReadFileRecords(path string, fn RecordReaderFunc) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()
	reader := bufio.NewReader(file)
	return ReadRecords(reader, fn)
}

func ReadRecords(reader *bufio.Reader, fn RecordReaderFunc) error {
	jsonBytes, err := reader.ReadBytes(RS)
	if err != nil {
		return ErrNotJsonSeq
	}
	res := bytes.TrimSpace(jsonBytes)
	if len(res) > 1 {
		return ErrExtraneousBytes
	}
	for {
		jsonBytes, err := reader.ReadBytes(RS)
		if err == nil {
			err = trimBytes(jsonBytes[:len(jsonBytes)-1], fn)
			if err != nil {
				return err
			}
		} else if err == io.EOF {
			err = fn(bytes.TrimSpace(jsonBytes), err)
			if err != nil {
				return err
			}
			return io.EOF
		} else {
			return err
		}
	}
}

func trimBytes(b []byte, fn RecordReaderFunc) error {
	res := bytes.TrimSpace(b)
	if len(res) > 0 {
		return fn(res, nil)
	} else {
		return fn(res, ErrEmptyPayload)
	}
}
