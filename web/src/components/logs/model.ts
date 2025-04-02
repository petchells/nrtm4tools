export interface LogLine {
	time: string;
	level: string;
	msg: string;
	[p: string]: any;
}

export interface UserMessage {
	ID: string;
	Content: LogLine;
}

export enum ToolbarCommand {
	closeLogPane = "CloseLogPane",
	reconnectWS = "ReconnectWebSocket",
}

// type IsPrimitive<T> = keyof T extends never ? true : false

// type DeepReadonly<T> = {
//   readonly [P in keyof T]: IsPrimitive<T[P]> extends true
//     ? T[P]
//     : DeepReadonly<T[P]>
// }

export function printParams(line: LogLine): string {
	const pmsg: string[] = [];
	for (let p in line) {
		if (["time", "level", "msg"].indexOf(p) === -1) {
			pmsg.push(`${p}=${line[p]}`)
		}
	}
	return pmsg.join(", ");
}
