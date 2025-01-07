import { useEffect, useState } from "react";
import Alert from "@mui/material/Alert";
import Box from "@mui/material/Box";
import CircularProgress from "@mui/material/CircularProgress";
import Grid from "@mui/material/Grid2";
import Typography from "@mui/material/Typography";
import WarningIcon from "@mui/icons-material/Warning";

import SourcesTable from "./SourcesTable.tsx";
import SourcesError from "./SourcesError.tsx";
import { SourceModel } from "../../client/models.ts";
import SourcesInput from "./SourcesInput.tsx";
import { WebAPIClient } from "../../client/WebAPIClient.ts";
import Source from "./Source.tsx";

export default function Sources() {
  const [pageLoading, setPageLoading] = useState<number>(1);
  const [loading, setLoading] = useState<number>(0);
  const [err, setErr] = useState<string>("");
  const [sources, setSources] = useState<SourceModel[]>([]);
  const [selectedIDs, setSelectedIDs] = useState<string[]>([]);
  const [refresh, setRefresh] = useState<number>(0);
  const client = new WebAPIClient();

  useEffect(() => {
    setPageLoading(1);
    if (sources.length) {
      console.log("why are we getting them again?");
    }
    client
      .listSources()
      .then(
        (ss) => {
          setSources(ss);
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
  }, []);

  const handleOnSelected = (row: SourceModel) => {
    const key = row.ID;
    const idx = selectedIDs.indexOf(key);
    if (idx < 0) {
      selectedIDs.push(key);
    } else {
      selectedIDs.splice(idx, 1);
    }
    setSelectedIDs(selectedIDs);
    setRefresh(refresh ^ 1);
  };

  const noSources = () => {
    return (
      <Alert
        icon={<WarningIcon fontSize="inherit" />}
        severity="success"
        sx={{ width: "100%", maxWidth: { sm: "100%", md: "1700px" } }}
      >
        No sources are available. You can add one with the 'connect' command.
      </Alert>
    );
  };

  const onUrlEntered = (url: string) => {
    console.log("url", url);
    setLoading(1);
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

  return (
    <Box sx={{ width: "100%", maxWidth: { sm: "100%", md: "1700px" } }}>
      <Typography variant="h4" component="h1" sx={{ mb: 2 }}>
        Sources
      </Typography>
      {!!pageLoading ? (
        <Box sx={{ display: "flex" }}>
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
                ></Source>
              ))}
          {!err && selectedIDs.length === 0 && (
            <SourcesInput onUrlEntered={onUrlEntered} disabled={!!loading} />
          )}
        </Grid>
      )}
    </Box>
  );
}
