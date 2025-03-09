import { useEffect, useState } from "react";

import Alert from "@mui/material/Alert";
import Box from "@mui/material/Box";
import Button from "@mui/material/Button";
import CircularProgress from "@mui/material/CircularProgress";
import Grid from "@mui/material/Grid2";
import Typography from "@mui/material/Typography";

import RefreshIcon from "@mui/icons-material/Refresh";
import WarningIcon from "@mui/icons-material/Warning";

import { SourceModel } from "../../client/models.ts";
import WebAPIClient from "../../client/WebAPIClient.ts";
import SourcesTable from "./SourcesTable.tsx";
import SourcesError from "./SourcesError.tsx";
import SourcesInput from "./SourcesInput.tsx";
import Source from "./Source.tsx";

export default function Sources() {
  const [pageLoading, setPageLoading] = useState<number>(1);
  const [dataLoading, setDataLoading] = useState(false);
  const [err, setErr] = useState<string>("");
  const [sources, setSources] = useState<SourceModel[]>([]);
  const [selectedIDs, setSelectedIDs] = useState<string[]>([]);
  const [refresh, setRefresh] = useState<number>(0);
  const client = new WebAPIClient();

  const fetchSources = () => {
    setPageLoading(1);
    client
      .listSources()
      .then(
        (ss) => {
          setSources(ss);
          for (let i = 0; i < selectedIDs.length; i++) {
            const found = ss.filter((s) => s.ID === selectedIDs[i]).length > 0;
            if (!found) {
              selectedIDs.splice(i, 1);
            }
          }
          setErr("");
        },
        (err) => {
          setSources([]);
          setSelectedIDs([]);
          setErr("Connection error: " + err);
        }
      )
      .then(() => {
        setPageLoading(0);
      });
  };

  useEffect(() => {
    fetchSources();
  }, []);

  const handleOnSelected = (row: SourceModel) => {
    const key = row.ID;
    const idx = selectedIDs.indexOf(key);
    if (idx < 0) {
      selectedIDs.push(key);
    } else {
      selectedIDs.splice(idx, 1);
    }
    const srcIDs = sources.map((ss) => ss.ID);
    selectedIDs.sort((a, b) => srcIDs.indexOf(a) - srcIDs.indexOf(b));
    setSelectedIDs(selectedIDs);
    setRefresh(refresh ^ 1);
  };

  const noSources = () => {
    return (
      <Alert
        icon={<WarningIcon fontSize="inherit" />}
        severity="info"
        sx={{ width: "100%", maxWidth: { sm: "100%", md: "1700px" } }}
      >
        No sources are available. Add one with the 'connect' command.
      </Alert>
    );
  };

  const onUrlEntered = (url: string, label: string) => {
    setDataLoading(true);
    client
      .connectSource(url, label)
      .then(
        (msg) => {
          console.log("success", msg);
        },
        (rej) => {
          if (typeof rej == "string") {
            setErr(rej);
          } else if (typeof rej == "object" && rej.message) {
            setErr(rej.message);
          } else {
            setErr("" + rej);
          }
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

  const handleSourceUpdated = (id: string, source: SourceModel) => {
    for (const s of sources) {
      if (s.ID === id) {
        s.Label = source.Label;
        setRefresh(refresh ^ 1);
        break;
      }
    }
  };

  const removeSource = (source: SourceModel) => {
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
        (err) => setErr(err)
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
    setErr("Source was not removed");
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
          {err ? (
            <SourcesError error={err} />
          ) : sources.length < 1 ? (
            noSources()
          ) : (
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
          {!err && selectedIDs.length === 0 && (
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
