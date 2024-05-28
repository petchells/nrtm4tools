package jsonseq

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"os"
	"strings"
)

// RS delimeter for records in a jsonseq
var RS byte = 0x1E

// ErrEmptyPayload returned when the payload is empty
var ErrEmptyPayload = errors.New("empty payload")

// ErrNotJSONSeq if the RS delimiter isn't found
var ErrNotJSONSeq = errors.New("not a JSON seq file")

// ErrExtraneousBytes returned when non-JSON chars are found in the payload
var ErrExtraneousBytes = errors.New("bytes found before record marker")

// RecordReaderFunc defines the callback function for jsonseq reads
type RecordReaderFunc func([]byte, error) error

// ReadStringRecords calls fn each time it finds a jsonseq record in jsonSeq
func ReadStringRecords(jsonSeq string, fn RecordReaderFunc) error {
	reader := bufio.NewReader(strings.NewReader(jsonSeq))
	return ReadRecords(reader, fn)
}

// ReadFileRecords reads a jsonseq file from path and calls fn for each record
func ReadFileRecords(path string, fn RecordReaderFunc) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()
	reader := bufio.NewReader(file)
	return ReadRecords(reader, fn)
}

// ReadRecords reads a jsonseq file and calls fn for each record
func ReadRecords(reader *bufio.Reader, fn RecordReaderFunc) error {
	jsonBytes, err := reader.ReadBytes(RS)
	if err != nil {
		return ErrNotJSONSeq
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
