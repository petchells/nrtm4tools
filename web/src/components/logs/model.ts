
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

// type IsPrimitive<T> = keyof T extends never ? true : false

// type DeepReadonly<T> = {
//   readonly [P in keyof T]: IsPrimitive<T[P]> extends true
//     ? T[P]
//     : DeepReadonly<T[P]>
// }

export function printLogLine(line: LogLine): string {
	let pmsg = line.msg;
	for (let p in line) {
		if (["time", "level", "msg"].indexOf(p) === -1) {
			console.log("xxxxx", p);
			pmsg+=`, ${p}=${line[p]}`
		}
	}
	return pmsg;
}
