import { useState } from "react";

import ClearIcon from "@mui/icons-material/Clear";
import SaveIcon from "@mui/icons-material/Save";
import RestartAltIcon from "@mui/icons-material/RestartAlt";

import OutlinedInput from "@mui/material/OutlinedInput";
import PunchClockIcon from "@mui/icons-material/PunchClock";
import Stack from "@mui/material/Stack";
import { IconButton } from "@mui/material";
import { formatDateWithStyle } from "../../util/dates";

export default function LabelControl(props: {
  value: string;
  disabled?: boolean;
  onTextEntered: (text: string) => void;
}) {
  const disableInput = !!props.disabled;
  const [inputText, setInputText] = useState(props.value);

  const initial = props.value;

  const btnClicked = () => {
    if (isValid(inputText)) {
      props.onTextEntered(inputText);
    }
  };

  const isValid = (str: string): boolean => {
    return str !== initial;
  };

  const clearInput = () => {
    setInputText("");
  };

  const resetInput = () => {
    setInputText(props.value);
  };

  const setLabelToTimestamp = () => {
    const ts = formatDateWithStyle(new Date(), "en-gb", "longdatetime");
    setInputText(ts);
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
        onClick={setLabelToTimestamp}
        disabled={disableInput || inputText.length > 0}
      >
        <PunchClockIcon />
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
