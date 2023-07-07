package nrtm4model

type FileRef struct {
	Version uint   `json:"version"`
	Url     string `json:"url"`
	Hash    string `json:"hash"`
}

type NrtmFile struct {
	NrtmVersion uint   `json:"nrtm_version"`
	Type        string `json:"type"`
	Source      string `json:"source"`
	SessionID   string `json:"session_id"`
	Version     uint   `json:"version"`
}

type Change struct {
	Action      string  `json:"action"`
	Object      *string `json:"object"`
	ObjectClass *string `json:"object_class"`
	PrimaryKey  *string `json:"primary_key"`
}

type Notification struct {
	NrtmFile
	Timestamp      string     `json:"timestamp"`
	NextSigningKey *string    `json:"next_signing_key"`
	Snapshot       FileRef    `json:"snapshot"`
	Deltas         *[]FileRef `json:"deltas"`
}

type DeltaFile struct {
	NrtmFile
}

type SnapshotFile struct {
	NrtmFile
}

type SnapshotObject struct {
	Object string `json:"object"`
}
