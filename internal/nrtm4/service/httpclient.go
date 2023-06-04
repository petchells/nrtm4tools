package service

import (
	"encoding/json"
	"io"
	"net/http"

	"gitlab.com/etchells/nrtm4client/internal/nrtm4/nrtm4model"
)

type httpClient struct{}

func (cl httpClient) getUpdateNotification(url string) (nrtm4model.Notification, error) {
	var file nrtm4model.Notification
	if err := cl.fetchObject(url, &file); err != nil {
		return file, err
	}
	return file, nil
}

func (cl httpClient) fetchFile(url string) (io.ReadCloser, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	return resp.Body, err
}

func (cl httpClient) fetchObject(url string, file any) error {
	var resp *http.Response
	var err error
	if resp, err = http.Get(url); err != nil {
		return err
	}
	if err = json.NewDecoder(resp.Body).Decode(&file); err != nil {
		return err
	}
	return nil
}
