import Box from "@mui/material/Box";
import Grid from "@mui/material/Grid2";
import Typography from "@mui/material/Typography";

export default function Logs() {
  return (
    <Box sx={{ width: "100%", maxWidth: { sm: "100%", md: "1700px" } }}>
      <Grid container spacing={2} columns={12}>
        <Grid>
          <Typography variant="body1" sx={{ mb: 2 }}>
            Coming soon
          </Typography>
        </Grid>
      </Grid>
    </Box>
  );
}
