import { ReactElement } from "react";

import Alert, { AlertColor } from "@mui/material/Alert";
import IconButton from "@mui/material/IconButton";
import CloseIcon from "@mui/icons-material/Close";
import ErrorIcon from "@mui/icons-material/Error";
import InfoIcon from "@mui/icons-material/Info";
import SuccessIcon from "@mui/icons-material/Check";
import WarningIcon from "@mui/icons-material/Warning";

interface AlertMessageProps {
  message: string;
  level?: AlertColor;
  dismissed: () => void;
}

export default function AlertMessage({
  message,
  level,
  dismissed,
}: AlertMessageProps) {
  let lvl: AlertColor = level || "success";
  let icon: ReactElement;
  switch (level) {
    case "error":
      icon = <ErrorIcon fontSize="inherit" />;
      break;
    case "warning":
      icon = <WarningIcon fontSize="inherit" />;
      break;
    case "info":
      icon = <InfoIcon fontSize="inherit" />;
      break;
    default:
      icon = <SuccessIcon fontSize="inherit" />;
  }
  return (
    <Alert
      icon={icon}
      severity={lvl}
      sx={{ width: "100%", maxWidth: { sm: "100%", md: "1700px" } }}
      action={
        <IconButton
          aria-label="close"
          color="inherit"
          size="small"
          onClick={dismissed}
        >
          <CloseIcon fontSize="inherit" />
        </IconButton>
      }
    >
      {message}
    </Alert>
  );
}
