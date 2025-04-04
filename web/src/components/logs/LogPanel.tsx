import { Fragment } from "react/jsx-runtime";

import CircleIcon from "@mui/icons-material/Circle";
import { createTheme, ThemeProvider, Typography } from "@mui/material";
import Grid from "@mui/material/Grid2";

import { LogLine, printParams } from "./model";
import { ReactElement } from "react";

const theme = createTheme({
  typography: {
    subtitle1: {
      fontSize: 12,
    },
    body1: {
      fontFamily: "monospace",
      fontSize: 12,
      wordBreak: "break-all",
    },
  },
});

const levels = ["ERROR", "WARN", "INFO", "DEBUG"];

const circledLetter = (letter: string, c: string) => {
  return (
    <Typography
      fontSize="small"
      sx={{
        color: c,
        background: c,
        borderRadius: 4,
        display: "inline-block",
        margin: "1px",
        width: "1.6em",
      }}
    >
      {letter}
    </Typography>
  );
};

const discs: { [lvl: string]: ReactElement } = {
  ERROR: circledLetter("E", "#cc0000"),
  WARN: circledLetter("W", "#6633ff"),
  INFO: circledLetter("I", "#339999"),
  DEBUG: circledLetter("D", "#99ccff"),
};

interface LogPanelProps {
  messageHistory: LogLine[];
  level: number;
}

export default function LogPanel({ messageHistory, level }: LogPanelProps) {
  return (
    <ThemeProvider theme={theme}>
      <Grid container spacing={0} columns={12}>
        {messageHistory
          .filter((log) => levels.indexOf(log.level) <= level)
          .map((line, idx) => (
            <Fragment key={line.time + idx}>
              <Grid size={4}>
                <Typography variant="body1">
                  {discs[line.level]} {line.time}
                </Typography>
              </Grid>
              <Grid size={8}>
                <Typography variant="body1">
                  <b>{line.msg}</b>
                </Typography>
                <Typography variant="body1">{printParams(line)}</Typography>
              </Grid>
            </Fragment>
          ))}
      </Grid>
    </ThemeProvider>
  );
}
