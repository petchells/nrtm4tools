import { Fragment } from "react/jsx-runtime";

import { createTheme, SxProps, ThemeProvider } from "@mui/material";
import Box from "@mui/material/Box";
import Fab from "@mui/material/Fab";
import Grid from "@mui/material/Grid2";
import Zoom from "@mui/material/Zoom";

import Typography from "@mui/material/Typography";
import WatchIcon from "@mui/icons-material/Watch";
import WatchOffIcon from "@mui/icons-material/WatchOff";
import { LogLine, printParams } from "./model";
import { useState } from "react";

const theme = createTheme({
  typography: {
    body1: {
      fontFamily: "monospace",
      fontSize: 12,
      wordBreak: "break-all",
    },
  },
});

const fabStyle: SxProps = {
  position: "absolute",
  bottom: 16,
  right: 16,
};

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
  const [scrollBottom, setScrollBottom] = useState(true);
  let to;
  clearTimeout(to);
  to = setTimeout(() => {
    if (!scrollBottom) {
      return;
    }
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
      <Box sx={{ "& > :not(style)": { m: 1 } }}>
        <Fab
          sx={fabStyle}
          size="small"
          color="secondary"
          aria-label="add"
          onClick={() => setScrollBottom(!scrollBottom)}
        >
          {scrollBottom ? <WatchIcon /> : <WatchOffIcon />}
        </Fab>
      </Box>
      <Box id="gridend"></Box>
    </ThemeProvider>
  );
}
