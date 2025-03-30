import { StrictMode } from "react";
import { createRoot } from "react-dom/client";
import { RouterProvider } from "react-router-dom";
import { ThemeProvider } from "@mui/material/styles";
import { CssBaseline } from "@mui/material";
import theme from "./theme";
import router from "./root";
import { AppConfig } from "./client/models";

// fetch("/s/clientcfg.json")
//   .then(
//     (resp) => resp.json(),
//   )
//   .then(
//     (cfg: AppConfig) => {
//       console.log(cfg);
//     },
//     () => console.log("Failed to get configuration"),
//   );
export var websocketURL = await fetch("/s/clientcfg.json")
  .then((resp) => resp.json())
  .then(
    (cfg: AppConfig) => cfg.WebSocketURL,
    () => "/ws"
  );

createRoot(document.getElementById("root")!).render(
  <StrictMode>
    <ThemeProvider theme={theme}>
      <CssBaseline />
      <RouterProvider router={router} />
    </ThemeProvider>
  </StrictMode>
);
