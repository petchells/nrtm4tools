package service

import (
	"errors"
	"net/url"
	"strings"

	"github.com/petchells/nrtm4client/internal/nrtm4/persist"
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

type fileRefsByVersion []persist.FileRefJSON

func (s fileRefsByVersion) Len() int {
	return len(s)
}
func (s fileRefsByVersion) Less(i, j int) bool {
	return s[i].Version < s[j].Version
}

func (s fileRefsByVersion) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
