
//
// NRTM4 files
//

export interface FileRef {
	hash: string;
	url: string;
	version: number;
}

export interface NotificationJSON {
	deltas: FileRef[];
	next_signing_key: string;
	nrtm_version: number;
	session_id: string;
	snapshot: FileRef;
	source: string;
	timestamp: string;
	type: string;
	version: number;
}

//
// Server models
//

export interface Notification {
	ID: string;
	Created: string;
	SourceID: string;
	Payload: NotificationJSON;
	SessionID: string;
	Source: string;
	Version: number;
}


export interface SourceDetail {
	ID: string;
	Source: string;
	SessionID: string;
	Version: number;
	NotificationURL: string;
	Label: string;
	Status: string;
	Created: string;
	Properties: SourceProperties;
	Notifications: Notification[];
}

export interface SourceProperties {
	AutoUpdate: AutoUpdateMode;
	Compliance: ComplianceMode;
}

enum AutoUpdateMode {
	Off = 0,
	Preserve = 1,
	Replace = 2,
}

enum ComplianceMode {
	Loose = 0,
	Strict = 1,
}

export interface AppConfig {
	WebSocketURL: string;
	RPCEndpoint: string;
}
