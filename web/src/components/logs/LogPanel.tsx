import { Fragment } from "react/jsx-runtime";

import { Box, createTheme, ThemeProvider, Typography } from "@mui/material";
import Grid from "@mui/material/Grid2";

import { LogLine, printParams } from "./model";

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

const showDate = (dstr: string): string => {
  const str = dstr.substring(11, 23);
  return str;
};

interface LogPanelProps {
  messageHistory: LogLine[];
  level: number;
}

export default function LogPanel({ messageHistory, level }: LogPanelProps) {
  let to;
  clearTimeout(to);
  to = setTimeout(() => {
    const el = document.getElementById("gridend");
    el && el.scrollIntoView();
  }, 800);
  return (
    <ThemeProvider theme={theme}>
      <Grid container spacing={0} columns={12}>
        {messageHistory
          .filter((log) => levels.indexOf(log.level) <= level)
          .map((line) => (
            <Fragment key={line.time}>
              <Grid size={2} sx={{ maxWidth: 160 }}>
                <Typography variant="body1">
                  <span className={"loglevel " + line.level}>
                    {line.level[0]}
                  </span>
                  {showDate(line.time)}
                </Typography>
              </Grid>
              <Grid size={10}>
                <Typography variant="body1">
                  <b>{line.msg}</b> {printParams(line)}
                </Typography>
              </Grid>
            </Fragment>
          ))}
      </Grid>
      <Box id="gridend"></Box>
    </ThemeProvider>
  );
}
