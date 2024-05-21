package service

import "testing"

func TestURLUtil(t *testing.T) {
	var fileName string
	{
		fileName, _ = fileNameFromURLString(stubNotificationURL)
		if fileName != "notification.json" {
			t.Error("file name is wrong")
		}
	}
	{
		fileName, _ = fileNameFromURLString("https://nrtm.db.ripe.net/nrtmv4/RIPE/update-notification-file.json?aparam=50")
		if fileName != "update-notification-file.json" {
			t.Error("file name should be", "update-notification-file.json")
		}
	}
}
