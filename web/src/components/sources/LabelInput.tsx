import { useState } from "react";

import Button from "@mui/material/Button";
import ButtonGroup from "@mui/material/ButtonGroup";
import OutlinedInput from "@mui/material/OutlinedInput";
import PunchClockIcon from "@mui/icons-material/PunchClock";
import Stack from "@mui/material/Stack";
import { IconButton } from "@mui/material";
import { formatDateWithStyle } from "../../util/dates";

export default function LabelInput(props: {
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

  const setLabelToTimestamp = () => {
    const ts = formatDateWithStyle(new Date(), "longdatetime");
    setInputText(ts);
  };

  return (
    <Stack direction="row" spacing={2}>
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
        onClick={setLabelToTimestamp}
        disabled={disableInput || inputText.length > 0}
      >
        <PunchClockIcon />
      </IconButton>
      <ButtonGroup>
        <Button
          onClick={clearInput}
          disabled={disableInput || inputText.length === 0}
        >
          Clear
        </Button>
        <Button
          onClick={btnClicked}
          disabled={disableInput || !isValid(inputText)}
        >
          Save
        </Button>
      </ButtonGroup>
    </Stack>
  );
}
