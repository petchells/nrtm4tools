import { expect, test } from "vitest";
import { formatDateWithStyle, parseISOString } from "./dates";

test("parseISOString() should convert an ISO datetime to a native date with ms resolution", () => {
  {
    const dateStr = "2023-10-10T15:15:15.987654321Z"; // Golang RFC3339 resolves to nanoseconds, PostgreSQL to microseconds
    const dt = parseISOString(dateStr);
    expect(dt).toBeTruthy();
    expect(dt.getFullYear()).toBe(2023);
    expect(dt.getUTCMonth()).toBe(9);
    expect(dt.getUTCDate()).toBe(10);
    expect(dt.getUTCHours()).toBe(15);
    expect(dt.getUTCMinutes()).toBe(15);
    expect(dt.getUTCSeconds()).toBe(15);
    expect(dt.getUTCMilliseconds()).toBe(987);
  }
  {
    const dateStr = "2023-10-10T15:15:15.9";
    const dt = parseISOString(dateStr);
    expect(dt.getUTCMilliseconds()).toBe(900);
  }
  {
    const dateStr = "2023-10-10T15:15:15.09";
    const dt = parseISOString(dateStr);
    expect(dt.getUTCMilliseconds()).toBe(90);
  }
  {
    const dateStr = "2023-10-10T15:15:15.009";
    const dt = parseISOString(dateStr);
    expect(dt.getUTCMilliseconds()).toBe(9);
  }
});

test("parseISOString() should convert an ISO date to a native date", () => {
  const dateStr = "2023-10-10";
  const dt = parseISOString(dateStr);
  expect(dt).toBeTruthy();
  expect(dt.getFullYear()).toBe(2023);
  expect(dt.getUTCMonth()).toBe(9);
  expect(dt.getUTCDate()).toBe(10);
  expect(dt.getUTCHours()).toBe(0);
  expect(dt.getUTCMinutes()).toBe(0);
});

test("formatDateWithStyle() from a date string works for a full datetime", () => {
  const dateStr = "2023-10-10T15:15:15.987654321Z";
  const result = formatDateWithStyle(dateStr, "en-gb", "longdatetime");
  const expected = " October 2023 at "; // the rest depends on the tz of the machine running the test
  expect(result).toContain(expected);
});
