import { useState } from "react";

import ClearIcon from "@mui/icons-material/Clear";
import IconButton from "@mui/material/IconButton";
import OutlinedInput from "@mui/material/OutlinedInput";
import RestartAltIcon from "@mui/icons-material/RestartAlt";
import SaveIcon from "@mui/icons-material/Save";
import Stack from "@mui/material/Stack";

interface LabelControlProps {
  value: string;
  disabled?: boolean;
  onTextEntered: (text: string) => void;
}

export default function LabelControl({
  value,
  disabled,
  onTextEntered,
}: LabelControlProps) {
  const disableInput = !!disabled;
  const [inputText, setInputText] = useState(value);

  const initial = value;

  const btnClicked = () => {
    if (isValid(inputText)) {
      onTextEntered(inputText);
    }
  };

  const isValid = (str: string): boolean => {
    return str !== initial;
  };

  const clearInput = () => {
    setInputText("");
  };

  const resetInput = () => {
    setInputText(value);
  };

  return (
    <Stack direction="row" spacing={1}>
      <OutlinedInput
        disabled={disableInput}
        size="small"
        id="labelinput"
        value={inputText}
        sx={{ flexGrow: 1 }}
        onChange={(event: React.ChangeEvent<HTMLInputElement>) => {
          setInputText(event.target.value);
        }}
        inputProps={{
          "aria-label": "Label text",
        }}
      />
      <IconButton
        onClick={clearInput}
        disabled={disableInput || inputText.length === 0}
      >
        <ClearIcon />
      </IconButton>
      <IconButton
        onClick={resetInput}
        disabled={disableInput || !isValid(inputText)}
      >
        <RestartAltIcon />
      </IconButton>
      <IconButton
        onClick={btnClicked}
        disabled={disableInput || !isValid(inputText)}
      >
        <SaveIcon />
      </IconButton>
    </Stack>
  );
}
