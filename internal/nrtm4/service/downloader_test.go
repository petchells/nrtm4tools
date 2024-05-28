package service

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"gitlab.com/etchells/nrtm4client/internal/nrtm4/nrtm4model"
	"gitlab.com/etchells/nrtm4client/internal/nrtm4/persist"
)

func TestSuccess(t *testing.T) {
	dl := downloader{}
	_, errs := dl.downloadNotificationFile(dlClientStub{}, "")
	if len(errs) > 0 {
		t.Error("should not be any errors but found:", errs[0])
	} else {
		t.Log("OK")
	}
}

type dlClientStub struct {
	Client
}

func (c dlClientStub) getUpdateNotification(string) (nrtm4model.NotificationJSON, error) {
	notification := nrtm4model.NotificationJSON{
		NrtmFileJSON: nrtm4model.NrtmFileJSON{
			NrtmVersion: 4,
			SessionID:   uuid.NewString(),
			Type:        persist.NotificationFile.String(),
			Source:      "ZZZZ",
			Version:     22,
		},
		Timestamp:      time.Now().UTC().Format(time.RFC3339),
		NextSigningKey: new(string),
		SnapshotRef: nrtm4model.FileRefJSON{
			URL: "https://xxx.xxx.xx/notification.json",
		},
		DeltaRefs: &[]nrtm4model.FileRefJSON{},
	}
	return notification, nil
}
