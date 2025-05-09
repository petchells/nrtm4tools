import { useState } from "react";
import { styled } from "@mui/material/styles";
import Box from "@mui/material/Box";
import FormControl from '@mui/material/FormControl';
import FormControlLabel from "@mui/material/FormControlLabel";
import IconButton from "@mui/material/IconButton";
// import Input from "@mui/material/Input";
import InputLabel from '@mui/material/InputLabel';
import MenuItem from '@mui/material/MenuItem';
import Paper from "@mui/material/Paper";
import Radio from "@mui/material/Radio";
import RadioGroup from "@mui/material/RadioGroup";
import Select, { SelectChangeEvent } from '@mui/material/Select';
import Stack from "@mui/material/Stack";

import SaveIcon from "@mui/icons-material/Save";

import { UpdateMode, SourceProperties } from "../../client/models";
import { Typography } from "@mui/material";

const PaddedItem = styled(Paper)(({ theme }) => ({
    ...theme.typography.body2,
    padding: theme.spacing(2),
    textAlign: "start",
}));

interface SourcePropertiesEditorProperties {
    sourceProps: SourceProperties;
    saveSourceProps: (p: SourceProperties) => void;
}

export default function SourcePropertiesEditor({
    sourceProps,
    saveSourceProps,
}: SourcePropertiesEditorProperties) {

    const [updateMode, setUpdateMode] = useState(sourceProps.UpdateMode || UpdateMode.Preserve);
    const [autoUpdateInterval, setAutoUpdateInterval] = useState(sourceProps.AutoUpdateInterval || 0);


    const handleIntervalChange = (event: SelectChangeEvent) => {
        const val = parseInt(event.target.value as string, 10);
        setAutoUpdateInterval(val);
    };

    const handleUpdateChange = (event: SelectChangeEvent) => {
        const val = parseInt(event.target.value as string, 10);
        console.log(val)
        setUpdateMode(val);
    };

    const propertiesHaveChanged = () => {
        // Deliberate non-use of !== cz sourceProps may be empty, which is still valid
        return updateMode != sourceProps.UpdateMode || autoUpdateInterval != sourceProps.AutoUpdateInterval;
    };

    return (
        <PaddedItem>
            <Stack spacing={2}>
                <FormControl fullWidth>
                    <InputLabel id="poll-interval-label">Notification poll interval</InputLabel>
                    <Select
                        labelId="poll-interval-label"
                        id="poll-interval-select"
                        value={"" + autoUpdateInterval}
                        onChange={handleIntervalChange}
                    >
                        <MenuItem value={0}>Off</MenuItem>
                        <MenuItem value={1}>One minute</MenuItem>
                        <MenuItem value={2}>Two minutes</MenuItem>
                        <MenuItem value={5}>Five minutes</MenuItem>
                        <MenuItem value={15}>Fifteen minutes</MenuItem>
                        <MenuItem value={30}>Thirty minutes</MenuItem>
                        <MenuItem value={60}>One hour</MenuItem>
                    </Select>
                </FormControl>
                <RadioGroup>
                    <Typography>When reinitializing from a snapshot:</Typography>
                    <FormControlLabel
                        label="Preserve stale repositories"
                        control={<Radio
                            size="small"
                            checked={updateMode === UpdateMode.Preserve}
                            onChange={handleUpdateChange}
                            value={UpdateMode.Preserve}
                            name="autoupdate-radio-button"
                            slotProps={{
                                input: {
                                    'aria-label': 'Preserve broken repository when reinitializing',
                                },
                            }}
                        />} />
                    <FormControlLabel
                        label="Delete old repository"
                        control={<Radio
                            size="small"
                            checked={updateMode === UpdateMode.Replace}
                            onChange={handleUpdateChange}
                            value={UpdateMode.Replace}
                            name="autoupdate-radio-button"
                            slotProps={{
                                input: {
                                    'aria-label': 'Replace repository when reinitializing',
                                },
                            }}
                        />} />
                </RadioGroup>
            </Stack>
            <Box sx={{ mt: 1, width: "100%" }}>
                <IconButton onClick={() => saveSourceProps({ AutoUpdateInterval: autoUpdateInterval, UpdateMode: updateMode })} disabled={!propertiesHaveChanged()}>
                    <SaveIcon />
                </IconButton>
            </Box>
        </PaddedItem>
    );
}