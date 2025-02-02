import Alert from "@mui/material/Alert";
import ErrorIcon from "@mui/icons-material/Error";

export default function SourcesError(props: { error: string }) {
  const err = props.error;
  const errorContainer = (err: string) => {
    return (
      <Alert
        icon={<ErrorIcon fontSize="inherit" />}
        severity="error"
        sx={{ width: "100%", maxWidth: { sm: "100%", md: "1700px" } }}
      >
        Error: {err}
      </Alert>
    );
  };
  return errorContainer(err);
}
