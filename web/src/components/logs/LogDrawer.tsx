import { useState, useEffect, useCallback } from "react";
import useWebSocket, * as ws from "react-use-websocket";

import Box from "@mui/material/Box";
import Drawer from "@mui/material/Drawer";

import Stack from "@mui/material/Stack";

import FrameToolbar from "./FrameToolbar";
import LogPanel from "./LogPanel";
import { websocketURL } from "../../main";
import { LogLine, UserMessage } from "./model";

const drawerHeight = 360;

interface LogDrawerProps {
  open: boolean;
  setOpen: (b: boolean) => void;
}

export default function LogDrawer({ open, setOpen }: LogDrawerProps) {
  const [messageHistory, setMessageHistory] = useState<LogLine[]>([]);
  const { sendMessage, lastMessage, readyState } = useWebSocket(websocketURL);
  //const [refresh, setRefresh] = useState<boolean>(false);

  useEffect(() => {
    if (lastMessage !== null) {
      const msg: UserMessage = JSON.parse(lastMessage.data);
      //setRefresh(!refresh);
      setMessageHistory((prev) => prev.concat(msg.Content));
    }
  }, [lastMessage]);

  const msg = {
    ID: "logs",
    Content: "Hello",
  };

  const handleClickSendMessage = useCallback(
    () => sendMessage(JSON.stringify(msg)),
    []
  );

  const connectionStatus = {
    [ws.ReadyState.CONNECTING]: "Connecting",
    [ws.ReadyState.OPEN]: "Open",
    [ws.ReadyState.CLOSING]: "Closing",
    [ws.ReadyState.CLOSED]: "Closed",
    [ws.ReadyState.UNINSTANTIATED]: "Uninstantiated",
  }[readyState];

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
        <FrameToolbar setOpen={setOpen} status={connectionStatus} />
        <Stack>
          <Box sx={{ m: 2 }}>
            <LogPanel messageHistory={messageHistory} />
          </Box>
        </Stack>
      </Box>
    </Drawer>
  );
}
