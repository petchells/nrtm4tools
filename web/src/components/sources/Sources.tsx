import { ReactElement, useEffect, useState } from "react";

import Alert, { AlertColor } from "@mui/material/Alert";
import Box from "@mui/material/Box";
import Button from "@mui/material/Button";
import CircularProgress from "@mui/material/CircularProgress";
import Grid from "@mui/material/Grid2";
import Typography from "@mui/material/Typography";

import ErrorIcon from "@mui/icons-material/Error";
import RefreshIcon from "@mui/icons-material/Refresh";

import { SourceDetail } from "../../client/models.ts";
import WebAPIClient from "../../client/WebAPIClient.ts";
import SourcesTable from "./SourcesTable.tsx";
import SourcesInput from "./SourcesInput.tsx";
import Source from "./Source.tsx";
import { JsonRPCError } from "../../client/RPCClient.ts";

export default function Sources() {
  const [pageLoading, setPageLoading] = useState(false);
  const [dataLoading, setDataLoading] = useState(false);
  const [alert, setAlert] = useState<ReactElement | null>(null);
  const [sources, setSources] = useState<SourceDetail[]>([]);
  const [selectedIDs, setSelectedIDs] = useState<string[]>([]);
  const [refresh, setRefresh] = useState<number>(0);
  const client = new WebAPIClient();

  const bySourceThenLabel = (a: SourceDetail, b: SourceDetail) => {
    if (a.Source === b.Source) {
      return a.Label.localeCompare(b.Label);
    }
    return a.Source.localeCompare(b.Source);
  };
  const newMessage = (msg: string | JsonRPCError, level?: AlertColor) => {
    let icon;
    if (!level) {
      level = "info";
      icon = <ErrorIcon fontSize="inherit" />;
    }
    let text: string;
    if (typeof msg === "string") {
      text = msg;
    } else {
      text = `RPCError${msg.code} ${msg.message}`;
    }
    setAlert(
      <Alert
        icon={icon}
        severity={level}
        sx={{ width: "100%", maxWidth: { sm: "100%", md: "1700px" } }}
      >
        {text}
      </Alert>
    );
  };

  const fetchSources = () => {
    setPageLoading(true);
    client
      .listSources()
      .then(
        (ss) => {
          setSources(ss.sort(bySourceThenLabel));
          const ids = ss.map((s) => s.ID);
          const sIDs = selectedIDs.map((sid) => sid);
          for (let i = 0; i < sIDs.length; i++) {
            if (ids.indexOf(sIDs[i]) < 0) {
              selectedIDs.splice(selectedIDs.indexOf(sIDs[i]), 1);
            }
          }
          if (ss.length === 0) {
            newMessage(
              "No sources are available. Add one with the 'connect' command.",
              "info"
            );
          } else {
            setAlert(null);
          }
        },
        (err) => {
          setSources([]);
          setSelectedIDs([]);
          if (err.hasOwnProperty("message")) {
            newMessage("Error: " + err.message, "error");
          } else {
            newMessage("Connection error: " + err, "error");
          }
        }
      )
      .finally(() => {
        setPageLoading(false);
      });
  };

  let timeout: number;

  useEffect(() => {
    clearTimeout(timeout);
    timeout = setTimeout(() => fetchSources(), 50);
  }, []);

  const handleOnSelected = (row: SourceDetail) => {
    const key = row.ID;
    const idx = selectedIDs.indexOf(key);
    if (idx > -1) {
      setSelectedIDs([]);
    } else {
      setSelectedIDs([key]);
    }
    //    setRefresh(refresh ^ 1);
  };

  const onUrlEntered = (url: string, label: string) => {
    setDataLoading(true);
    client
      .connectSource(url, label)
      .then(
        (msg) => {
          console.log("success", msg);
        },
        (err) => {
          newMessage(err, "error");
        }
      )
      .finally(() => setDataLoading(false));
  };

  const lookupSource = (key: string) => {
    const src = sources.filter((s) => {
      return key === s.ID;
    });
    return src[0];
  };

  const handleSourceUpdated = (id: string, source: SourceDetail) => {
    for (let i = 0; i < sources.length; i++) {
      if (sources[i].ID === id) {
        const src = sources[i];
        src.Label = source.Label;
        src.Notifications.splice(
          0,
          src.Notifications.length,
          ...source.Notifications
        );
        src.Version = source.Version;
        src.Status = source.Status;
        setRefresh(refresh ^ 1);
        break;
      }
    }
  };

  const removeSource = (source: SourceDetail) => {
    setDataLoading(true);
    client
      .removeSource(source.Source, source.Label)
      .then(
        () => {
          const idx = selectedIDs.indexOf(source.ID);
          if (idx > -1) {
            selectedIDs.splice(idx, 1);
            setSelectedIDs(selectedIDs);
          }
          const oidx = sources.indexOf(source);
          if (oidx > -1) {
            sources.splice(oidx, 1);
            setSources(sources);
          }
          setRefresh(refresh ^ 1);
        },
        (err) => newMessage(err, "error")
      )
      .finally(() => setDataLoading(false));
  };

  const handleRemoveSource = (id: string) => {
    for (const s of sources) {
      if (s.ID === id) {
        removeSource(s);
        return;
      }
    }
    newMessage("Source was not removed", "error");
  };

  return (
    <Box sx={{ width: "100%", maxWidth: { sm: "100%", md: "1700px" } }}>
      <Grid container spacing={2} columns={12}>
        <Box sx={{ mb: 1 }}>
          <Button
            variant="outlined"
            size="small"
            onClick={fetchSources}
            startIcon={<RefreshIcon />}
          >
            Refresh
          </Button>
        </Box>
      </Grid>
      {!!pageLoading ? (
        <Box sx={{ display: "flex", justifyContent: "center", mt: 5 }}>
          <CircularProgress />
        </Box>
      ) : (
        <Grid container spacing={2} columns={12}>
          {alert}
          {sources.length > 0 && (
            <SourcesTable
              rows={sources}
              selectedIDs={selectedIDs}
              onSelected={handleOnSelected}
            />
          )}
          {selectedIDs.length > 0 &&
            selectedIDs
              .map((key) => lookupSource(key))
              .map((src) => (
                <Source
                  key={src.ID}
                  source={src}
                  sourceUpdated={handleSourceUpdated}
                  sourceRemoveConfirmed={handleRemoveSource}
                ></Source>
              ))}
          {selectedIDs.length === 0 && (
            <>
              <Typography variant="h6" component="h1" sx={{ mt: 2 }}>
                Connect a new source
              </Typography>
              <SourcesInput
                onUrlEntered={onUrlEntered}
                disabled={!!dataLoading}
              />
            </>
          )}
        </Grid>
      )}
    </Box>
  );
}
