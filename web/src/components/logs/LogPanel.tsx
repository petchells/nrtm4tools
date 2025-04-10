import { Fragment } from "react/jsx-runtime";

import { createTheme, ThemeProvider } from "@mui/material";
import Box from "@mui/material/Box";
import Grid from "@mui/material/Grid2";
import Typography from "@mui/material/Typography";

import { LogLine, printParams } from "./model";

const theme = createTheme({
  typography: {
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
  scrollBottom: boolean;
}

let to: number;

export default function LogPanel({
  messageHistory,
  level,
  scrollBottom,
}: LogPanelProps) {
  clearTimeout(to);
  to = setInterval(() => {
    if (!scrollBottom) {
      return;
    }
    const el = document.getElementById("gridend");
    const rect = el?.getBoundingClientRect();
    const inView =
      rect &&
      rect.bottom <=
        (window.innerHeight || document.documentElement.clientHeight);
    !inView && el && el.scrollIntoView();
  }, 600);

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
