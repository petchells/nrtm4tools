package service

import (
	"encoding/json"
	"io"
	"net/http"

	"gitlab.com/etchells/nrtm4client/internal/nrtm4/nrtm4model"
)

// Client fetches things from the NRTM server, or anywhwere, actually
type Client interface {
	getUpdateNotification(string) (nrtm4model.NotificationJSON, error)
	getResponseBody(string) (io.Reader, error)
}

// HTTPClient implementation of Client
type HTTPClient struct{}

func (cl HTTPClient) getUpdateNotification(url string) (nrtm4model.NotificationJSON, error) {
	var file nrtm4model.NotificationJSON
	if err := cl.getObject(url, &file); err != nil {
		return file, err
	}
	return file, nil
}

func (cl HTTPClient) getResponseBody(url string) (io.Reader, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	return resp.Body, err
}

func (cl HTTPClient) getObject(url string, obj any) error {
	var resp *http.Response
	var err error
	if resp, err = http.Get(url); err != nil {
		return err
	}
	if err = json.NewDecoder(resp.Body).Decode(&obj); err != nil {
		return err
	}
	return nil
}
