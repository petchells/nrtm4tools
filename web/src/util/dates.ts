export function parseISOString(s: string) {
	const sp = s.split(/\D+/);
	const b = sp.map((s) => parseInt(s, 10)).filter((d) => !isNaN(d));
	let i = b.length;
	if (i < 3 || i > 7) {
		console.log("invalid date", s, b);
		throw Error("invalid date");
	}
	let ms = 0;
	if (sp.length > 6) {
		const msStr = sp[6].substring(0, 3); // round to ms
		switch (msStr.length) {
			case 3:
				ms = parseInt(msStr, 10);
				break;
			case 2:
				ms = parseInt(msStr, 10) * 10;
				break;
			case 1:
				ms = parseInt(msStr, 10) * 100;
				break;
			default:
		}
	}
	const b7: number[] = [0, 0, 0, 0, 0, 0, 0];
	for (let j = 0; j < i; j++) {
		if (j === 1) {
			b7[j] = --b[j];
		} else if (j == 6) {
			b7[j] = ms;
		} else {
			b7[j] = b[j];
		}
	}
	for (; i < 7; i++) {
		b7[i] = 0;
	}
	return new Date(Date.UTC(b7[0], b7[1], b7[2], b7[3], b7[4], b7[5], b7[6]));
}

const options: { [key: string]: any } = {
	longdatetime: {
		year: "numeric",
		month: "long",
		day: "2-digit",
		hour: "2-digit",
		minute: "2-digit",
		hour12: false,
	},
	long: {
		year: "numeric",
		month: "long",
		day: "2-digit",
	},
	short: {
		year: "numeric",
		month: "short",
		day: "2-digit",
	},
	narrow: {
		month: "2-digit",
		day: "2-digit",
	},
};

type styles = "longdatetime" | "long" | "short" | "narrow";

export function formatDateWithStyle(
	date: string | Date,
	locale: string,
	style?: styles
): string {
	if (!date) {
		return "";
	}
	try {
		const d = typeof date === "string" ? new Date(date) : date;
		return Intl.DateTimeFormat(locale, options[style || "short"]).format(d);
	} catch (e) {
		console.log("Invalid date string", date);
		return "";
	}
}
