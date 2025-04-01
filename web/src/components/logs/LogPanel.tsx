import Grid from "@mui/material/Grid2";
import { LogLine } from "./model";

interface LogPanelProps {
  messageHistory: LogLine[];
}

export default function LogPanel({ messageHistory }: LogPanelProps) {
  return (
    <Grid container spacing={0} columns={12}>
      {messageHistory.map((line) => (
        <>
          <Grid size={3}>{line.time}</Grid>
          <Grid size={1}>{line.level}</Grid>
          <Grid size={8}>{line.msg}</Grid>
        </>
      ))}
    </Grid>
  );
}
