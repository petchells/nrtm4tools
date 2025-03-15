import { useState, useEffect, useCallback } from "react";
import useWebSocket, * as ws from 'react-use-websocket';
import Box from "@mui/material/Box";
import Grid from "@mui/material/Grid2";

export default function Logs() {
  //Public API that will echo messages sent to it back to the client
  const [socketUrl, setSocketUrl] = useState('http://localhost:8090/ws');
  const [messageHistory, setMessageHistory] = useState<MessageEvent<any>[]>([]);

  const { sendMessage, lastMessage, readyState } = useWebSocket(socketUrl);

  useEffect(() => {
    if (lastMessage !== null) {
      setMessageHistory((prev) => prev.concat(lastMessage));
    }
  }, [lastMessage]);

  const handleClickChangeSocketUrl = useCallback(
    () => setSocketUrl('wss://demos.kaazing.com/echo'),
    []
  );

  const handleClickSendMessage = useCallback(() => sendMessage('Hello'), []);

  const connectionStatus = {
    [ws.ReadyState.CONNECTING]: 'Connecting',
    [ws.ReadyState.OPEN]: 'Open',
    [ws.ReadyState.CLOSING]: 'Closing',
    [ws.ReadyState.CLOSED]: 'Closed',
    [ws.ReadyState.UNINSTANTIATED]: 'Uninstantiated',
  }[readyState];

  return (
    <Box sx={{ width: "100%", maxWidth: { sm: "100%", md: "1700px" } }}>
      <Grid container spacing={2} columns={12}>
        <Grid size={12}>
          <button onClick={handleClickChangeSocketUrl}>
            Click Me to change Socket Url
          </button>
          <button
            onClick={handleClickSendMessage}
            disabled={readyState !== ws.ReadyState.OPEN}
          >
            Click Me to send 'Hello'
          </button>
        </Grid>
        <Grid size={12}>
          The WebSocket is currently {connectionStatus}
        </Grid>
        <Grid size={12}>
          {lastMessage ? <span>Last message: {lastMessage.data}</span> : null}
        </Grid>
        <Grid size={12}>
          {messageHistory.map((message, idx) => (
            <div key={idx}>{message ? message.data : null}</div>
          ))}
        </Grid>
      </Grid>
    </Box >
  );
};
