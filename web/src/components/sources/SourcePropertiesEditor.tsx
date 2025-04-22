import { styled } from "@mui/material/styles";
import Box from "@mui/material/Box";
// import Button from "@mui/material/Button";
import FormControlLabel from "@mui/material/FormControlLabel";
import IconButton from "@mui/material/IconButton";
import Input from "@mui/material/Input";
import Paper from "@mui/material/Paper";
// import Stack from "@mui/material/Stack";
import Radio from "@mui/material/Radio";
import RadioGroup from "@mui/material/RadioGroup";
//import Typography from "@mui/material/Typography";
import SaveIcon from "@mui/icons-material/Save";

import { UpdateMode, SourceProperties } from "../../client/models";
import { useState } from "react";

const PaddedItem = styled(Paper)(({ theme }) => ({
    ...theme.typography.body2,
    padding: theme.spacing(1),
    paddingLeft: theme.spacing(2),
    textAlign: "start",
}));

interface SourcePropertiesEditorProperties {
    sourceProps: SourceProperties;
}

export default function SourcePropertiesEditor({
    sourceProps,
}: SourcePropertiesEditorProperties) {

    const [props, setProps] = useState(sourceProps);

    return (
        <PaddedItem>
            <Box>
                <RadioGroup>
                    <FormControlLabel
                        label="Off"
                        value="off"
                        control={<Radio
                            size="small"
                            checked={!props.UpdateMode}
                            onChange={() => { }}
                            name="autoupdate-radio-button"
                            slotProps={{
                                input: {
                                    'aria-label': 'Autoupdate off',
                                },
                            }}
                        />} />
                    <FormControlLabel
                        label="Preserve"
                        value="preserve"
                        control={<Radio
                            size="small"
                            checked={props.UpdateMode === UpdateMode.Preserve}
                            onChange={() => { }}
                            value="preserve"
                            name="autoupdate-radio-button"
                            slotProps={{
                                input: {
                                    'aria-label': 'Preserve broken repository when reinitializing',
                                },
                            }}
                        />} />
                    <FormControlLabel
                        label="Replace"
                        value="replace"
                        control={<Radio
                            size="small"
                            checked={props.UpdateMode === UpdateMode.Replace}
                            onChange={() => { }}
                            value="replace"
                            name="autoupdate-radio-button"
                            slotProps={{
                                input: {
                                    'aria-label': 'Replace repository when reinitializing',
                                },
                            }}
                        />} />
                </RadioGroup>
                <FormControlLabel
                    label="Interval (minutes)"
                    control={<Input />}
                />
            </Box>
            <Box sx={{ mt: 1 }}>
                <IconButton>
                    <SaveIcon />
                </IconButton>
            </Box>
        </PaddedItem>
    );
}