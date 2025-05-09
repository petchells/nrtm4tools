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

const levels: string[] = ["ERROR", "WARN", "INFO", "DEBUG"];
const levelInitials = "!•⁃‣";

const showDate = (dstr: string): string => {
  const str = dstr.substring(11, 23);
  return str;
};

interface LogPanelProps {
  messageHistory: LogLine[];
  level: number;
  scrollBottom: boolean;
}

const scrollElementIntoView = (id: string) => {
  const el = document.getElementById(id);
  if (!el) {
    return;
  }
  const rect = el.getBoundingClientRect();
  const inView =
    rect.bottom <=
    (window.innerHeight || document.documentElement.clientHeight);
  !inView && el.scrollIntoView();
};

let to: number;

export default function LogPanel({
  messageHistory,
  level,
  scrollBottom,
}: LogPanelProps) {
  clearTimeout(to);
  if (scrollBottom) {
    to = setInterval(() => scrollElementIntoView("gridend"), 600);
  }

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
                    {levelInitials[levels.indexOf(line.level)]}
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
