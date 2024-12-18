import { useRouteError } from "react-router-dom";
import type {} from "@mui/x-date-pickers/themeAugmentation";
import type {} from "@mui/x-charts/themeAugmentation";
import type {} from "@mui/x-data-grid/themeAugmentation";
import type {} from "@mui/x-tree-view/themeAugmentation";
import { ThemeProvider } from "@mui/material/styles";
import Box from "@mui/material/Box";
import Stack from "@mui/material/Stack";
import CssBaseline from "@mui/material/CssBaseline";
import { StrictMode } from "react";
import AppTheme from "./components/shared-theme/AppTheme";
import theme from "./theme";

import {
  chartsCustomizations,
  dataGridCustomizations,
  datePickersCustomizations,
  treeViewCustomizations,
} from "./theme/customizations";

const xThemeComponents = {
  ...chartsCustomizations,
  ...dataGridCustomizations,
  ...datePickersCustomizations,
  ...treeViewCustomizations,
};

export default function ErrorPage(props: { disableCustomTheme?: boolean }) {
  const error = useRouteError() as { statusText: string; message: string };
  console.error(error);

  return (
    <StrictMode>
      <ThemeProvider theme={theme}>
        <CssBaseline />

        <AppTheme {...props} themeComponents={xThemeComponents}>
          <CssBaseline enableColorScheme />

          <Box sx={{ display: "flex" }}>
            <Box
              component="main"
              sx={(theme) => ({
                flexGrow: 1,
                backgroundColor: theme.vars
                  ? `rgba(${theme.vars.palette.background.defaultChannel} / 1)`
                  : alpha(theme.palette.background.default, 1),
                overflow: "auto",
              })}
            >
              <Stack
                spacing={2}
                sx={{
                  alignItems: "center",
                  mx: 3,
                  pb: 5,
                  mt: { xs: 8, md: 0 },
                }}
              >
                <div id="error-page">
                  <h1>Oops!</h1>
                  <p>Sorry, an unexpected error has occurred.</p>
                  <p>
                    <i>{error.statusText || error.message}</i>
                  </p>
                </div>
              </Stack>
            </Box>
          </Box>
        </AppTheme>
      </ThemeProvider>
    </StrictMode>
  );
}
