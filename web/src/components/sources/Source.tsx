import { useState } from "react";

import { styled } from "@mui/material/styles";
import Box from "@mui/material/Box";
import Button from "@mui/material/Button";
import CircularProgress from "@mui/material/CircularProgress";
import Dialog from "@mui/material/Dialog";
import DialogTitle from "@mui/material/DialogTitle";
import DialogContent from "@mui/material/DialogContent";
import DialogContentText from "@mui/material/DialogContentText";
import DialogActions from "@mui/material/DialogActions";
import Grid from "@mui/material/Grid2";
import Link from "@mui/material/Link";
import Paper from "@mui/material/Paper";
import Stack from "@mui/material/Stack";
import Typography from "@mui/material/Typography";

import DeleteOutlineIcon from '@mui/icons-material/DeleteOutline';
import UpdateIcon from '@mui/icons-material/Update';

import { formatDateWithStyle, parseISOString } from "../../util/dates";
import { SourceModel } from "../../client/models";
import WebAPIClient from "../../client/WebAPIClient.ts";
import LabelControl from "./LabelControl.tsx";

export default function Source(props: {
  source: SourceModel;
  sourceUpdated: (id: string, source: SourceModel) => void;
  sourceRemoveConfirmed: (sourceID: string) => void;
}) {
  const client = new WebAPIClient();
  const source = props.source;

  const [loading, setLoading] = useState(false);
  const [open, setOpen] = useState(false);

  const removeSourceClicked = () => {
    setOpen(true);
  };

  const updateSourceClicked = () => {
    setLoading(true);
    client
      .updateSource(source.Source, source.Label)
      .then(() => {
        props.sourceUpdated(source.ID, source);
      }, (msg) => console.log(msg))
      .finally(() => setLoading(false));
  };

  const saveLabel = (text: string) => {
    setLoading(true);
    client
      .saveLabel(source.Source, source.Label, text)
      .then(
        (resp) => {
          source.Label = resp.Label;
          props.sourceUpdated(source.ID, source);
        },
        (err) => console.log(err)
      )
      .finally(() => setLoading(false));
  };

  const handleClose = (confirm: boolean) => () => {
    setOpen(false);
    if (confirm) {
      props.sourceRemoveConfirmed(source.ID);
    }
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
    <Box sx={{ width: "100%", maxWidth: { sm: "100%", md: "1700px" }, mt: 4 }}>
      <Dialog
        open={open}
        onClose={handleClose}
        aria-labelledby="alert-dialog-title"
        aria-describedby="alert-dialog-description"
      >
        <DialogTitle id="alert-dialog-title">
          {"Confirm removal of source"}
        </DialogTitle>
        <DialogContent>
          <DialogContentText id="alert-dialog-description">
            This will remove {source.Source} {source.Label} and all its objects and history from the repository.<br></br>
            Click 'Confirm' to continue.
          </DialogContentText>
        </DialogContent>
        <DialogActions>
          <Button onClick={handleClose(false)}>Cancel</Button>
          <Button onClick={handleClose(true)} autoFocus>Confirm</Button>
        </DialogActions>
      </Dialog>

      <Typography variant="h4" component="h2" sx={{ mb: 2 }}>
        {source.Source} {source.Label}
      </Typography>
      <Stack spacing={1} direction="row" sx={{ mb: 1 }}>
        <Button
          variant="outlined"
          size="small"
          startIcon={<UpdateIcon />}
          onClick={updateSourceClicked}
        >
          Update
        </Button>
        <Button
          variant="outlined"
          color="error"
          size="small"
          startIcon={<DeleteOutlineIcon />}
          onClick={removeSourceClicked}>
          Remove
        </Button>
      </Stack>
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
            <LabelControl
              value={source.Label}
              onTextEntered={saveLabel}
            ></LabelControl>
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
            <Link href={source.NotificationURL} target="_blank" rel="noopener">
              {source.NotificationURL}
              <sup>🔗</sup>
            </Link>
          </Item>
        </Grid>
        <Grid size={{ xs: 4, md: 4 }}>
          <Label>Repo last updated</Label>
        </Grid>
        <Grid size={{ xs: 8, md: 8 }}>
          <Item>
            {source.Notifications.length > 0 ?
              formatDateWithStyle(
                parseISOString(source.Notifications[0].Created),
                "en-gb",
                "longdatetime"
              ) : (
                <i>Not synced</i>
              )}
          </Item>
        </Grid>
      </Grid>
    </Box>
  );
}
