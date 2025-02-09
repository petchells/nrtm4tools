package service

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/petchells/nrtm4client/internal/nrtm4/persist"
)

// HTTPResponseError is used to model an error response from a http client
type HTTPResponseError struct {
	Message string
	Status  int
	URL     string
}

func (cerr HTTPResponseError) Error() string {
	return fmt.Sprintln("HTTPClientError", cerr.URL, cerr.Status, cerr.Message)
}

// Client fetches things from the NRTM server, or anywhwere, actually
type Client interface {
	getUpdateNotification(string) (persist.NotificationJSON, error)
	getResponseBody(string) (io.Reader, error)
}

// HTTPClient implementation of Client
type HTTPClient struct{}

func (cl HTTPClient) getUpdateNotification(url string) (persist.NotificationJSON, error) {
	var unf persist.NotificationJSON
	body, err := cl.getResponseBody(url)
	if err != nil {
		return unf, err
	}
	err = json.NewDecoder(body).Decode(&unf)
	return unf, err
}

func (cl HTTPClient) getResponseBody(url string) (io.Reader, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode == http.StatusOK {
		return resp.Body, err
	}
	logger.Warn("HTTPClient getResponseBody received bad response", "status", resp.StatusCode, "message", resp.Status)
	return nil, clientErrFromResponse(resp)
}

func clientErrFromResponse(resp *http.Response) HTTPResponseError {
	return HTTPResponseError{Status: resp.StatusCode, Message: resp.Status, URL: resp.Request.URL.String()}
}
