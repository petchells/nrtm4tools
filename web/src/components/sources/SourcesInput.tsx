import { useState } from "react";

import IconButton from "@mui/material/IconButton";
import InputAdornment from "@mui/material/InputAdornment";
import NetworkPingIcon from "@mui/icons-material/NetworkPing";
import OutlinedInput from "@mui/material/OutlinedInput";
import PlayArrowIcon from "@mui/icons-material/PlayArrow";
import { grey } from "@mui/material/colors";

export default function SourcesInput(props: {
  disabled: boolean;
  onUrlEntered: (url: string) => void;
}) {
  const disableInput = props.disabled;
  const [inputText, setInputText] = useState("");

  const btnClicked = () => {
    if (isValidHttpUrl(inputText)) {
      props.onUrlEntered(inputText);
    }
  };

  const isValidHttpUrl = (str: string): boolean => {
    let url;
    try {
      url = new URL(str);
    } catch (_) {
      return false;
    }
    return url.protocol === "http:" || url.protocol === "https:";
  };

  return (
    <OutlinedInput
      disabled={disableInput}
      size="small"
      id="nurl"
      placeholder="Notification URL…"
      sx={{ flexGrow: 1 }}
      onChange={(event: React.ChangeEvent<HTMLInputElement>) => {
        setInputText(event.target.value);
      }}
      startAdornment={
        <InputAdornment position="start" sx={{ color: "text.primary" }}>
          <NetworkPingIcon
            sx={disableInput ? { color: grey[400] } : { color: grey[800] }}
            fontSize="small"
          />
        </InputAdornment>
      }
      endAdornment={
        <InputAdornment position="end" sx={{ color: "text.primary" }}>
          <IconButton
            aria-label="connect"
            onClick={btnClicked}
            edge="end"
            disabled={disableInput || !isValidHttpUrl(inputText)}
          >
            <PlayArrowIcon />
          </IconButton>
        </InputAdornment>
      }
      inputProps={{
        "aria-label": "search",
      }}
    />
  );
}
