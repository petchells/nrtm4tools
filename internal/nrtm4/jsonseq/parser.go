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

type UnmarshalFunc func([]byte, error)

func ParseString(jsonSeq string, fn UnmarshalFunc) error {
	reader := bufio.NewReader(strings.NewReader(jsonSeq))
	return ParseReader(reader, fn)
}

func ParseFile(path string, fn UnmarshalFunc) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()
	reader := bufio.NewReader(file)
	return ParseReader(reader, fn)
}

func ParseReader(reader *bufio.Reader, fn UnmarshalFunc) error {
	for {
		jsonBytes, err := reader.ReadBytes(RS)
		if err == nil {
			trimBytes(jsonBytes[:len(jsonBytes)-1], fn)
		} else if err == io.EOF {
			fn(bytes.TrimSpace(jsonBytes), err)
			return io.EOF
		} else {
			return err
		}
	}
}

func trimBytes(b []byte, fn UnmarshalFunc) {
	res := bytes.TrimSpace(b)
	if len(res) > 0 {
		fn(res, nil)
	} else {
		fn(res, ErrEmptyPayload)
	}
}
