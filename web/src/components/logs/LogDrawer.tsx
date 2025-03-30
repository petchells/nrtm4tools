//import { styled } from "@mui/material/styles";
import Box from "@mui/material/Box";
import Drawer from "@mui/material/Drawer";
//import MuiToolbar from "@mui/material/Toolbar";
import Stack from "@mui/material/Stack";

import FrameToolbar from "./FrameToolbar";
import LogPanel from "./LogPanel";

// const Toolbar = styled(MuiToolbar)({
//   width: "100%",
//   padding: "12px",
//   display: "flex",
//   flexDirection: "column",
//   alignItems: "start",
//   justifyContent: "center",
//   gap: "12px",
//   flexShrink: 0,
// });

const drawerHeight = 360;

interface LogDrawerProps {
  open: boolean;
  setOpen: (b: boolean) => void;
}

export default function LogDrawer({ open, setOpen }: LogDrawerProps) {
  const toolClick = () => {
    console.log("clicked tool");
    setOpen(false);
  };

  return (
    <Drawer
      sx={{
        height: drawerHeight,
        flexShrink: 0,
        "& .MuiDrawer-paper": {
          height: drawerHeight,
        },
      }}
      variant="persistent"
      anchor="bottom"
      open={open}
    >
      <Box sx={{ width: "100%", maxWidth: { sm: "100%", md: "1700px" } }}>
        <Stack>
          <FrameToolbar setOpen={setOpen} />
          <Box sx={{ m: 2 }}>
            <LogPanel />
          </Box>
        </Stack>
      </Box>
    </Drawer>
  );
}
