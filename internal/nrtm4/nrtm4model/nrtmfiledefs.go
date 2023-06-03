package nrtm4model

type FileRef struct {
	Version uint   `json:"version"`
	Url     string `json:"url"`
	Hash    string `json:"hash"`
}

type Change struct {
	Action      string  `json:"action"`
	Object      *string `json:"object"`
	ObjectClass *string `json:"object_class"`
	PrimaryKey  *string `json:"primary_key"`
}

type Notification struct {
	NrtmVersion int        `json:"nrtm_version"`
	Timestamp   string     `json:"timestamp"`
	Type        string     `json:"type"`
	Source      string     `json:"source"`
	SessionID   string     `json:"session_id"`
	Version     int        `json:"version"`
	Snapshot    FileRef    `json:"snapshot"`
	Deltas      *[]FileRef `json:"deltas"`
}

type DeltaFile struct {
	NrtmVersion uint     `json:"nrtm_version"`
	Type        string   `json:"type"`
	Source      string   `json:"source"`
	SessionID   string   `json:"session_id"`
	Version     uint     `json:"version"`
	Changes     []Change `json:"changes"`
}

type SnapshotFile struct {
	NrtmVersion uint     `json:"nrtm_version"`
	Type        string   `json:"type"`
	Source      string   `json:"source"`
	SessionID   string   `json:"session_id"`
	Version     uint     `json:"version"`
	Objects     []string `json:"objects"`
}
