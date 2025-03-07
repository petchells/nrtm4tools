import Box from "@mui/material/Box";
import Grid from "@mui/material/Grid2";
import Typography from "@mui/material/Typography";

export default function Queries() {
  return (
    <Box sx={{ width: "100%", maxWidth: { sm: "100%", md: "1700px" } }}>
      <Typography variant="h3" component="h1" sx={{ mb: 2 }}>
        Object queries
      </Typography>
      <Grid container spacing={2} columns={12}>
        <Grid>
          <Typography variant="body1" sx={{ mb: 2 }}>
            Coming soon-ish.
          </Typography>
        </Grid>
      </Grid>
    </Box>
  );
}
