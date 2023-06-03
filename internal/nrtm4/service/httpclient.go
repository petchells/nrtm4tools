package service

import (
	"encoding/json"
	"net/http"

	"gitlab.com/etchells/nrtm4client/internal/nrtm4/nrtm4model"
)

type nrtmHttpClient struct {
	url string
}

func newNrtmHttpClient(url string) nrtmHttpClient {
	return nrtmHttpClient{url}
}

func (c nrtmHttpClient) getUpdateNotification() (nrtm4model.Notification, error) {
	var file nrtm4model.Notification
	var resp *http.Response
	var err error
	if resp, err = http.Get(c.url); err != nil {
		return file, err
	}
	if err = json.NewDecoder(resp.Body).Decode(&file); err != nil {
		return file, err
	}
	return file, nil
}
