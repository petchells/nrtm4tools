import { useEffect, useState } from "react";
import Alert from "@mui/material/Alert";
import Box from "@mui/material/Box";
import Button from "@mui/material/Button";
import ButtonGroup from "@mui/material/ButtonGroup";
import CircularProgress from "@mui/material/CircularProgress";
import Grid from "@mui/material/Grid2";
import Typography from "@mui/material/Typography";
import WarningIcon from "@mui/icons-material/Warning";

import RpcClientService from "../../client/rpcClientService.ts";
import SourcesTable from "./SourcesTable.tsx";
import SourcesError from "./SourcesError.tsx";
import { SourceModel } from "../../client/models.ts";
import SourcesInput from "./SourcesInput.tsx";

export default function Sources() {
  const [pageLoading, setPageLoading] = useState<number>(1);
  const [loading, setLoading] = useState<number>(0);
  const [err, setErr] = useState<string>("");
  const [sources, setSources] = useState<SourceModel[]>([]);
  const [selectedIDs, setSelectedIDs] = useState<string[]>([]);
  const [refresh, setRefresh] = useState<number>(0);

  useEffect(() => {
    const rpcService = new RpcClientService();
    setPageLoading(1);
    rpcService
      .execute<SourceModel[]>("ListSources")
      .then(
        (ss) => {
          setSources(ss);
          setErr("");
        },
        () => {
          setSources([]);
          setSelectedIDs([]);
          setErr("No connection to the server");
        }
      )
      .then(() => {
        setPageLoading(0);
      });
  }, []);

  const handleOnSelected = (row: SourceModel) => {
    const key = row.Source + "." + row.Label;
    const idx = selectedIDs.indexOf(key);
    if (idx < 0) {
      selectedIDs.push(key);
    } else {
      selectedIDs.splice(idx, 1);
    }
    console.log(selectedIDs);
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
          {selectedIDs.length > 0 && (
            <ButtonGroup variant="outlined" aria-label="Actions for source">
              <Button>Label</Button>
              <Button>Update</Button>
            </ButtonGroup>
          )}
          {!err && selectedIDs.length === 0 && (
            <SourcesInput onUrlEntered={onUrlEntered} disabled={!!loading} />
          )}
        </Grid>
      )}
    </Box>
  );
}
