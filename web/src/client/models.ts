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

export interface Notification {
	ID: string;
	Created: string;
	SourceID: string;
	Payload: NotificationJSON;
	SessionID: string;
	Source: string;
	Version: number;
}

export interface SourceModel {
	ID: string;
	Source: string;
	SessionID: string;
	Version: number;
	NotificationURL: string;
	Label: string;
	Status: string;
	Created: string;
	Notifications: Notification[];
}

export interface AppConfig {
	WebSocketURL: string;
	RPCEndpoint: string;
}
