package persist

import "testing"

func TestFileTypeString(t *testing.T) {
	{
		expected := "notification"
		ft1 := NTRMFileType(NotificationFile)
		if ft1.String() != expected {
			t.Error("File type should be", expected)
		}
	}
	{
		expected := "snapshot"
		ft1 := NTRMFileType(SnapshotFile)
		if ft1.String() != expected {
			t.Error("File type should be", expected)
		}
	}
	{
		expected := "delta"
		ft1 := NTRMFileType(DeltaFile)
		if ft1.String() != expected {
			t.Error("File type should be", expected)
		}
	}
	{
		expected := ""
		ft1 := NTRMFileType(666)
		if ft1.String() != expected {
			t.Error("File type should be", expected)
		}
	}
}

func TestStringFileType(t *testing.T) {
	{
		ft, err := ToFileType("notification")
		if err != nil {
			t.Error("Unexpected error", err)
		}
		if ft != NotificationFile {
			t.Error("Expected", NotificationFile, "but was", ft)
		}
	}
	{
		ft, err := ToFileType("delta")
		if err != nil {
			t.Error("Unexpected error", err)
		}
		if ft != DeltaFile {
			t.Error("Expected", NotificationFile, "but was", ft)
		}
	}
	{
		ft, err := ToFileType("snapshot")
		if err != nil {
			t.Error("Unexpected error", err)
		}
		if ft != SnapshotFile {
			t.Error("Expected", NotificationFile, "but was", ft)
		}
	}
	{
		_, err := ToFileType("")
		if err == nil {
			t.Error("Expected an error")
		}
	}
	{
		_, err := ToFileType("nosuchfiletype")
		if err == nil {
			t.Error("Expected an error")
		}
	}
}
