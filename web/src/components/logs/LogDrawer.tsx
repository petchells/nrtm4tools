import { useState, useEffect } from "react";
import useWebSocket, * as ws from "react-use-websocket";

import Box from "@mui/material/Box";
import Drawer from "@mui/material/Drawer";

import Stack from "@mui/material/Stack";

import FrameToolbar from "./FrameToolbar";
import LogPanel from "./LogPanel";
import { websocketURL } from "../../main";
import { LogLine, ToolbarCommand, UserMessage } from "./model";

const drawerHeight = 360;

interface LogDrawerProps {
  open: boolean;
  setOpen: (b: boolean) => void;
}

export default function LogDrawer({ open, setOpen }: LogDrawerProps) {
  const [messageHistory, setMessageHistory] = useState<LogLine[]>([]);
  const [wsURL, setWsURL] = useState(websocketURL);
  const { lastMessage, readyState } = useWebSocket(wsURL);

  useEffect(() => {
    if (lastMessage !== null) {
      try {
        const msg: UserMessage = JSON.parse(lastMessage.data);
        setMessageHistory((prev) => prev.concat(msg.Content));
      } catch (ex) {
        console.log("lastMessage", lastMessage, ex);
      }
    }
  }, [lastMessage]);

  //   const msg = {
  //     ID: "logs",
  //     Content: "Hello",
  //   };

  const connectionStatus = {
    [ws.ReadyState.CONNECTING]: "Connecting",
    [ws.ReadyState.OPEN]: "Open",
    [ws.ReadyState.CLOSING]: "Closing",
    [ws.ReadyState.CLOSED]: "Closed",
    [ws.ReadyState.UNINSTANTIATED]: "Uninstantiated",
  }[readyState];

  const reconnect = () => {
    setWsURL("");
    setTimeout(() => setWsURL(websocketURL), 1000);
  };

  const logLevelChanged = (lvl: string) => {
    console.log("level", lvl);
  };
  const toolbarClick = (cmd: ToolbarCommand, ...args: any) => {
    switch (cmd) {
      case ToolbarCommand.closeLogPane:
        setOpen(false);
        return;
      case ToolbarCommand.reconnectWS:
        reconnect();
        return;
      case ToolbarCommand.setLogLevel:
        logLevelChanged(args[0]);
        return;
      default:
        throw "Not a ToolbarCommand";
    }
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
      <Box sx={{ width: "100%", maxWidth: { sm: "100%", md: "xl" } }}>
        <FrameToolbar toolbarClick={toolbarClick} status={connectionStatus} />
        <Stack>
          <Box sx={{ m: 2 }}>
            <LogPanel messageHistory={messageHistory} />
          </Box>
        </Stack>
      </Box>
    </Drawer>
  );
}
