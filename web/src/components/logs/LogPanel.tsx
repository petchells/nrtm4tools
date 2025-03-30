import { useState, useEffect, useCallback } from "react";
import useWebSocket, * as ws from "react-use-websocket";
import Grid from "@mui/material/Grid2";
import { websocketURL } from "../../main";

export default function LogPanel() {
  const [messageHistory, setMessageHistory] = useState<MessageEvent<any>[]>([]);
  const { sendMessage, lastMessage, readyState } = useWebSocket(websocketURL);

  useEffect(() => {
    if (lastMessage !== null) {
      setMessageHistory((prev) => prev.concat(lastMessage));
    }
  }, [lastMessage]);

  const msg = {
    ID: "public",
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
    <Grid container spacing={2} columns={12}>
      <Grid size={12}>
        <button
          onClick={handleClickSendMessage}
          disabled={readyState !== ws.ReadyState.OPEN}
        >
          Click Me to send 'Hello'
        </button>
      </Grid>
      <Grid size={12}>The WebSocket is currently {connectionStatus}</Grid>
      <Grid size={12}>
        {lastMessage ? <span>Last message: {lastMessage.data}</span> : null}
      </Grid>
      <Grid size={12}>
        {messageHistory.map((message, idx) => (
          <div key={idx}>{message ? message.data : null}</div>
        ))}
      </Grid>
    </Grid>
  );
}
