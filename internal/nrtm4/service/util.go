package service

import (
	"errors"
	"net/url"
	"strings"
)

func fileNameFromURLString(rawURL string) (string, error) {
	url, err := url.Parse(rawURL)
	if err != nil {
		return "", err
	}
	idx := strings.LastIndex(url.Path, "/")
	if idx > -1 {
		return url.Path[idx+1:], nil
	}
	return "", errors.New("did not find file name in URL")
}
