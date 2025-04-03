import { Fragment } from "react/jsx-runtime";

import CircleIcon from "@mui/icons-material/Circle";
import { createTheme, ThemeProvider, Typography } from "@mui/material";
import Grid from "@mui/material/Grid2";

import { LogLine, printParams } from "./model";
import { ReactElement } from "react";

interface LogPanelProps {
  messageHistory: LogLine[];
}

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

const discs: { [lvl: string]: ReactElement } = {
  ERROR: <CircleIcon fontSize="small" sx={{ color: "#cc0000" }} />,
  WARN: <CircleIcon fontSize="small" sx={{ color: "#6633ff" }} />,
  INFO: <CircleIcon fontSize="small" sx={{ color: "#339999" }} />,
  DEBUG: <CircleIcon fontSize="small" sx={{ color: "#99ccff" }} />,
};

export default function LogPanel({ messageHistory }: LogPanelProps) {
  return (
    <ThemeProvider theme={theme}>
      <Grid container spacing={0} columns={12}>
        {messageHistory.map((line) => (
          <Fragment key={line.time}>
            <Grid size={3}>
              <Typography variant="body1">{line.time}</Typography>
            </Grid>
            <Grid size={1}>
              <Typography variant="body1">{discs[line.level]}</Typography>
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
