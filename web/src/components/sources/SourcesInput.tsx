import { useState } from "react";

import { grey } from "@mui/material/colors";
import InputAdornment from "@mui/material/InputAdornment";
import NetworkPingIcon from "@mui/icons-material/NetworkPing";
import OutlinedInput from "@mui/material/OutlinedInput";
import Stack from "@mui/material/Stack";
import Button from "@mui/material/Button";

export default function SourcesInput(props: {
  disabled: boolean;
  onUrlEntered: (url: string, label: string) => void;
}) {
  const disableInput = props.disabled;
  const [inputText, setInputText] = useState("");
  const [labelText, setLabelText] = useState("");

  const btnClicked = () => {
    if (isValidHttpUrl(inputText)) {
      props.onUrlEntered(inputText, labelText);
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
    <Stack gap={1} direction={{ xs: "column", sm: "row" }} width="100%">
      <OutlinedInput
        id="nurlinput"
        placeholder="Notification URL…"
        disabled={disableInput}
        sx={{ flexGrow: 8 }}
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
        inputProps={{
          "aria-label": "URL input",
        }}
      />
      <OutlinedInput
        id="labelinput"
        placeholder="Label"
        disabled={disableInput}
        sx={{ flexGrow: 1 }}
        onChange={(event: React.ChangeEvent<HTMLInputElement>) => {
          setLabelText(event.target.value);
        }}
      />
      <Button
        variant="outlined"
        aria-label="connect"
        onClick={btnClicked}
        disabled={disableInput || !isValidHttpUrl(inputText)}
      >
        Connect
      </Button>
    </Stack>
  );
}
