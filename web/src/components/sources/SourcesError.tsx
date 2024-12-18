import Alert from "@mui/material/Alert";
import WarningIcon from "@mui/icons-material/Warning";

export default function SourcesError(props: { error: string }) {
  const err = props.error;
  const errorContainer = (err: string) => {
    return (
      <Alert
        icon={<WarningIcon fontSize="inherit" />}
        severity="warning"
        sx={{ width: "100%", maxWidth: { sm: "100%", md: "1700px" } }}
      >
        Cannot connect to backend: {err}
      </Alert>
    );
  };
  return errorContainer(err);
}
