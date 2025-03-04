package service

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/petchells/nrtm4tools/internal/nrtm4/testresources"
)

func TestGetNotification(t *testing.T) {
	var buf *strings.Builder
	var err error

	f := testresources.OpenFile(t, "update-notification-file.jose")
	buf = new(strings.Builder)
	_, err = io.Copy(buf, f)
	if err != nil {
		t.Fatal("Unexpected error reading unf")
	}
	f.Close()

	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "%v", buf)
	}))
	defer svr.Close()

	c := HTTPClient{}

	res, err := c.getUpdateNotification(svr.URL)
	if err != nil {
		t.Errorf("expected err to be nil got %v", err)
	}
	if err = validateNotificationFile(res); err != nil {
		t.Fatal("Notification file failed to validated:", err)
	}
	if res.Version != 399659 {
		t.Errorf("expected version to be %v got %v", 399659, res.Version)
	}
}

func TestGetResponseBody(t *testing.T) {
	expected := "dummy data"
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "%v", expected)
	}))
	defer svr.Close()

	c := HTTPClient{}
	res, err := c.getResponseBody(svr.URL)

	if err != nil {
		t.Errorf("expected err to be nil got %v", err)
	}
	buf := new(strings.Builder)
	_, err = io.Copy(buf, res)
	if err != nil {
		t.Fatal("Unexpected error")
	}
	if buf.String() != expected {
		t.Errorf("Expected buf to contain %v got %v", expected, res)
	}
}

func TestGetResponseBodyError(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	}))
	defer svr.Close()

	c := HTTPClient{}
	_, err := c.getResponseBody(svr.URL)

	rerr, ok := err.(HTTPResponseError)
	if !ok {
		t.Errorf("expected err to be %T but was %T", HTTPResponseError{}, err)
	}
	if rerr.Status != http.StatusForbidden {
		t.Fatal("Expected status to be", http.StatusForbidden, "but was", rerr.Status)
	}
}
