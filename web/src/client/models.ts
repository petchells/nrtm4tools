export type FileRef = {
  hash: string;
  url: string;
  version: number;
};

export type NotificationJSON = {
  deltas: FileRef[];
  next_signing_key: string;
  nrtm_version: number;
  session_id: string;
  snapshot: FileRef;
  source: string;
  timestamp: string;
  type: string;
  version: number;
};

export type Notification = {
  ID: string;
  Created: string;
  NRTMSourceID: string;
  Payload: NotificationJSON;
  SessionID: string;
  Source: string;
  Version: number;
};

export type SourceModel = {
  ID: string;
  Source: string;
  SessionID: string;
  Version: number;
  NotificationURL: string;
  Label: string;
  Created: string;
  Notifications: Notification[];
};
