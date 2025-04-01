
export interface LogLine {
	time: string;
	level: string;
	msg: string;
	[p: string]: any;
}

export interface UserMessage {
	ID: string;
	Content: string;
}
