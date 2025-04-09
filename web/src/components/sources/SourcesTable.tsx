import { useState } from "react";

import CircularProgress from "@mui/material/CircularProgress";
import Grid from "@mui/material/Grid2";
import Table from "@mui/material/Table";
import TableBody from "@mui/material/TableBody";
import TableCell from "@mui/material/TableCell";
import TableContainer from "@mui/material/TableContainer";
import TableHead from "@mui/material/TableHead";
import TableRow from "@mui/material/TableRow";
import Paper from "@mui/material/Paper";
import Checkbox from "@mui/material/Checkbox";
import ErrorIcon from "@mui/icons-material/Error";

import { SourceDetail } from "../../client/models.ts";
import { formatDateWithStyle } from "../../util/dates.ts";
import { styled } from "@mui/material";

interface SourcesTableProps {
  rows: SourceDetail[];
  selectedIDs: string[];
  onSelected: (row: SourceDetail) => void;
}
export default function SourcesTable({
  rows,
  selectedIDs,
  onSelected,
}: SourcesTableProps) {
  const [refresh, setRefresh] = useState<number>(0);

  const handleClick = (row: SourceDetail) => {
    onSelected(row);
    setRefresh(refresh ^ 1);
  };

  const rowIcon = (status: string) => {
    if (status === "new") {
      return <CircularProgress size="1em" />;
    } else if (status !== "ok") {
      return <ErrorIcon sx={{ fontSize: "1em" }} color="error" />;
    }
  };

  const Tickbox = styled(Checkbox)({
    borderRadius: 12,
  });

  return (
    <Grid size={{ xs: 12, lg: 12 }}>
      <TableContainer component={Paper}>
        <Table aria-label="Table of NRTM Sources" size={"medium"}>
          <TableHead>
            <TableRow>
              <TableCell padding="checkbox"></TableCell>
              <TableCell component="th" scope="row" padding="normal">
                Source
              </TableCell>
              <TableCell component="th" scope="row" padding="normal">
                Label
              </TableCell>
              <TableCell component="th" scope="row" padding="normal">
                Last updated
              </TableCell>
              <TableCell
                align="right"
                component="th"
                scope="row"
                padding="normal"
              >
                Version
              </TableCell>
            </TableRow>
          </TableHead>
          <TableBody>
            {rows.map((row, index) => {
              const isItemSelected = selectedIDs.includes(row.ID);
              const labelId = `enhanced-table-checkbox-${index}`;
              return (
                <TableRow
                  key={row.ID}
                  hover
                  onClick={() => handleClick(row)}
                  role="checkbox"
                  tabIndex={-1}
                  aria-checked={isItemSelected}
                  selected={isItemSelected}
                  sx={{ cursor: "pointer" }}
                >
                  <TableCell padding="checkbox">
                    <Tickbox
                      color="primary"
                      checked={isItemSelected}
                      slotProps={{ input: { "aria-labelledby": labelId } }}
                    />
                  </TableCell>
                  <TableCell
                    component="td"
                    id={labelId}
                    scope="row"
                    padding="normal"
                  >
                    {row.Source} {rowIcon(row.Status)}
                  </TableCell>
                  <TableCell component="td" scope="row" padding="normal">
                    {row.Label}
                  </TableCell>
                  <TableCell>
                    {!!row.Notifications.length &&
                      formatDateWithStyle(
                        row.Notifications[0].Created,
                        "en-gb",
                        "longdatetime"
                      )}
                  </TableCell>
                  <TableCell
                    align="right"
                    component="td"
                    scope="row"
                    padding="normal"
                  >
                    {row.Version}
                  </TableCell>
                </TableRow>
              );
            })}
          </TableBody>
        </Table>
      </TableContainer>
    </Grid>
  );
}
