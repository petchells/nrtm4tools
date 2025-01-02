import { useEffect, useState } from "react";

import { styled } from "@mui/material/styles";
import Box from "@mui/material/Box";
import Button from "@mui/material/Button";
import ButtonGroup from "@mui/material/ButtonGroup";
import CircularProgress from "@mui/material/CircularProgress";
import Paper from "@mui/material/Paper";
import Grid from "@mui/material/Grid2";
import Typography from "@mui/material/Typography";

import { SourceModel } from "../../client/models";
import { WebAPIClient } from "../../client/WebAPIClient.ts";
import { formatDateWithStyle, parseISOString } from "../../util/dates";
import LabelInput from "./LabelInput";

export default function Source(props: { source: SourceModel }) {
  const client = new WebAPIClient();
  const source = props.source;

  const [loading, setLoading] = useState<number>(0);

  // useEffect(() => {
  // }, []);

  const saveLabel = (text: string) => {
    console.log("change label to", text);
    setLoading(1);
    client
      .saveLabel(source.Source, source.Label, text)
      .then(
        (resp) => {
          source.Label = resp.Label;
        },
        (err) => console.log(err)
      )
      .then(() => setLoading(0));
  };

  const Label = styled(Paper)(({ theme }) => ({
    ...theme.typography.body2,
    padding: theme.spacing(1),
    textAlign: "end",
    color: theme.palette.text.secondary,
    ...theme.applyStyles("dark", {
      backgroundColor: "#1A2027",
    }),
    ...theme.applyStyles("light", {
      backgroundColor: "#EAF0F7",
    }),
  }));

  const Item = styled(Paper)(({ theme }) => ({
    ...theme.typography.body2,
    padding: theme.spacing(1),
    textAlign: "start",
  }));

  return (
    <Box sx={{ width: "100%", maxWidth: { sm: "100%", md: "1700px" } }}>
      <Typography variant="h4" component="h2" sx={{ mb: 2 }}>
        {source.Source} {source.Label}
      </Typography>
      <Grid container spacing={2}>
        <Grid size={{ xs: 4, md: 4 }}>
          <Label>Source</Label>
        </Grid>
        <Grid size={{ xs: 8, md: 8 }}>
          <Item>{source.Source}</Item>
        </Grid>
        <Grid size={{ xs: 4, md: 4 }}>
          <Label>Label</Label>
        </Grid>
        <Grid size={{ xs: 8, md: 8 }}>
          {!!loading ? (
            <CircularProgress />
          ) : (
            <LabelInput
              value={source.Label}
              onTextEntered={saveLabel}
            ></LabelInput>
          )}
        </Grid>
        <Grid size={{ xs: 4, md: 4 }}>
          <Label>Version</Label>
        </Grid>
        <Grid size={{ xs: 8, md: 8 }}>
          <Item>{source.Version}</Item>
        </Grid>
        <Grid size={{ xs: 4, md: 4 }}>
          <Label>Notification URL</Label>
        </Grid>
        <Grid size={{ xs: 8, md: 8 }}>
          <Item>
            <a target="_blank" href={source.NotificationURL}>
              {source.NotificationURL}
            </a>
          </Item>
        </Grid>
        <Grid size={{ xs: 4, md: 4 }}>
          <Label>Repo last updated</Label>
        </Grid>
        <Grid size={{ xs: 8, md: 8 }}>
          <Item>
            {formatDateWithStyle(
              parseISOString(source.Notifications[0].Created),
              "longdatetime"
            )}
          </Item>
        </Grid>
        <Grid size={{ xs: 12, md: 12 }}>
          <ButtonGroup variant="outlined" aria-label="Actions for source">
            <Button>Label</Button>
            <Button>Update</Button>
          </ButtonGroup>
        </Grid>
      </Grid>
    </Box>
  );
}
