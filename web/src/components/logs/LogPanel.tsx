import { createTheme, ThemeProvider, Typography } from "@mui/material";
import Grid from "@mui/material/Grid2";

import { LogLine, printParams } from "./model";
import { Fragment } from "react/jsx-runtime";

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

export default function LogPanel({ messageHistory }: LogPanelProps) {
  return (
    <ThemeProvider theme={theme}>
      <Grid container spacing={0} columns={12}>
        {messageHistory.map((line) => (
          <Fragment key="{line.time}">
            <Grid size={3}>
              <Typography variant="body1">{line.time}</Typography>
            </Grid>
            <Grid size={1}>
              <Typography variant="body1">{line.level}</Typography>
            </Grid>
            <Grid size={8}>
              <Typography variant="body1">
                <b>{line.msg}</b> {printParams(line)}
              </Typography>
            </Grid>
          </Fragment>
        ))}
      </Grid>
    </ThemeProvider>
  );
}
