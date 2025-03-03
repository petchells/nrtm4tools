package service

import (
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/golang-jwt/jwt"
	"github.com/petchells/nrtm4tools/internal/nrtm4/persist"
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

func (cl HTTPClient) getUpdateNotification(urlStr string) (persist.NotificationJSON, error) {
	nURL, err := url.Parse(urlStr)
	if err != nil {
		logger.Warn("Failed to parse URL", "urlStr", urlStr)
	}
	havePublicKey := strings.HasSuffix(nURL.Host, "ripe.net")
	// TODO: when url is RIPE domain cache this from https://ftp.ripe.net/ripe/dbase/nrtmv4/nrtmv4_public_key.txt
	keyTxt := `-----BEGIN PUBLIC KEY-----
MFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAEOkzpjobirEcqoR6zLXnPkm4cCTEY
Xi2rLlCSXc5EZ3L3PycAdDmWQtGHD8GF++RqWgrdKv+9l+InalmiCGkpRQ==
-----END PUBLIC KEY-----`

	var unf persist.NotificationJSON
	body, err := cl.getResponseBody(urlStr)
	if err != nil || body == nil {
		logger.Warn("Failed to read response", "urlStr", urlStr, "error", err)
		return unf, err
	}
	bytes, err := io.ReadAll(body)
	if err != nil {
		logger.Warn("Failed to read body", "urlStr", urlStr, "body", body, "error", err)
		return unf, err
	}
	var pub any
	if havePublicKey {
		block, _ := pem.Decode([]byte(keyTxt))
		if block == nil || block.Type != "PUBLIC KEY" {
			logger.Warn("Failed to decode PEM block containing public key", "urlStr", urlStr, "body", body, "error", err)
			return unf, errors.New("failed to decode PEM block containing public key")
		}
		pub, err = x509.ParsePKIXPublicKey(block.Bytes)
		if err != nil {
			logger.Warn("Failed to parse public key", "urlStr", urlStr, "pub", pub, "error", err)
			return unf, err
		}
	}
	tokenString := string(bytes)
	claims := jwt.MapClaims{}
	_, err = jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (any, error) {
		return pub, nil
	})
	if err != nil {
		logger.Warn("Failed to parse with claims", "urlStr", urlStr, "error", err)
		return unf, err
	}
	// do something with decoded claims
	cljson, err := json.Marshal(claims)
	if err != nil {
		logger.Warn("Failed to marshal claims", "urlStr", urlStr, "error", err)
	}
	notification := new(persist.NotificationJSON)
	err = json.Unmarshal(cljson, notification)
	if err != nil {
		logger.Warn("Failed to unmarshal claims", "urlStr", urlStr, "error", err)
	}
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
