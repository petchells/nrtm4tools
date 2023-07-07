package jsonseq

import (
	"bufio"
	"bytes"
	"io"
	"os"
	"strings"
)

var RS byte = 0x1E

type UnmarshalFunc func([]byte)

func ParseString(jsonSeq string, fn UnmarshalFunc) error {
	reader := bufio.NewReader(strings.NewReader(jsonSeq))
	return parseReader(reader, fn)
}

func ParseFile(path string, fn UnmarshalFunc) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()
	reader := bufio.NewReader(file)
	return parseReader(reader, fn)
}

func parseReader(reader *bufio.Reader, fn UnmarshalFunc) error {
	for {
		jsonBytes, err := reader.ReadBytes(RS)
		if err == nil {
			trimBytes(jsonBytes[:len(jsonBytes)-1], fn)
		} else if err == io.EOF {
			trimBytes(jsonBytes, fn)
			return nil
		} else {
			return err
		}
	}
}

func trimBytes(b []byte, fn UnmarshalFunc) {
	res := bytes.TrimSpace(b)
	if len(res) > 0 {
		fn(res)
	}
}
