package persist

const (
	// DeltaDeleteAction NRTM4 code for a delete operation
	DeltaDeleteAction string = "delete"
	// DeltaAddModifyAction NRTM4 code for an addition or modification operation
	DeltaAddModifyAction = "add_modify"
)

// FileRefJSON json model of a file reference in a Notification file
type FileRefJSON struct {
	Version int64  `json:"version"`
	URL     string `json:"url"`
	Hash    string `json:"hash"`
}

// NrtmFileJSON json model of fields common to all NRTM4 files
type NrtmFileJSON struct {
	NrtmVersion int64  `json:"nrtm_version"`
	Type        string `json:"type"`
	Source      string `json:"source"`
	SessionID   string `json:"session_id"`
	Version     int64  `json:"version"`
}

// DeltaJSON json model of a change record in a DeltaFile
type DeltaJSON struct {
	Action      string  `json:"action"`
	Object      *string `json:"object"`
	ObjectClass *string `json:"object_class"`
	PrimaryKey  *string `json:"primary_key"`
}

// NotificationJSON json model of an NRTM4 notification
type NotificationJSON struct {
	NrtmFileJSON
	Timestamp      string        `json:"timestamp"`
	NextSigningKey *string       `json:"next_signing_key"`
	SnapshotRef    FileRefJSON   `json:"snapshot"`
	DeltaRefs      []FileRefJSON `json:"deltas"`
}

// DeltaFileJSON json model of an NRTM4 delta file
type DeltaFileJSON struct {
	NrtmFileJSON
}

// SnapshotFileJSON json model of an NRTM4 snapshot
type SnapshotFileJSON struct {
	NrtmFileJSON
}

// SnapshotObjectJSON json model of an object record in a SnapshotFile
type SnapshotObjectJSON struct {
	Object string `json:"object"`
}
