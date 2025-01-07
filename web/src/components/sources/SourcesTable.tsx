import Grid from "@mui/material/Grid2";
import Table from "@mui/material/Table";
import TableBody from "@mui/material/TableBody";
import TableCell from "@mui/material/TableCell";
import TableContainer from "@mui/material/TableContainer";
import TableHead from "@mui/material/TableHead";
import TableRow from "@mui/material/TableRow";
import Paper from "@mui/material/Paper";
import { SourceModel } from "../../client/models.ts";
import Checkbox from "@mui/material/Checkbox";
import { useState } from "react";
import { formatDateWithStyle } from "../../util/dates.ts";

export default function SourcesTable(props: {
  rows: SourceModel[];
  selectedIDs: string[];
  onSelected: (row: SourceModel) => void;
}) {
  const [refresh, setRefresh] = useState<number>(0);

  const rows = props.rows;
  const selectedIDs = props.selectedIDs;

  const handleClick = (row: SourceModel) => {
    props.onSelected(row);
    setRefresh(refresh ^ 1);
  };

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
                  key={index}
                  hover
                  onClick={() => handleClick(row)}
                  role="checkbox"
                  aria-checked={isItemSelected}
                  tabIndex={-1}
                  selected={isItemSelected}
                  sx={{ cursor: "pointer" }}
                >
                  <TableCell padding="checkbox">
                    <Checkbox
                      color="primary"
                      checked={isItemSelected}
                      inputProps={{
                        "aria-labelledby": labelId,
                      }}
                    />
                  </TableCell>
                  <TableCell
                    component="td"
                    id={labelId}
                    scope="row"
                    padding="normal"
                  >
                    {row.Source}
                  </TableCell>
                  <TableCell component="td" scope="row" padding="normal">
                    {row.Label}
                  </TableCell>
                  <TableCell>
                    {formatDateWithStyle(
                      row.Notifications[0].Created,
                      "en-gb",
                      "short"
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
