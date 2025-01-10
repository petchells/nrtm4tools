import { useState } from "react";
import { Outlet, useLocation, useNavigate } from "react-router-dom";

import { alpha } from "@mui/material/styles";
import Box from "@mui/material/Box";
import CssBaseline from "@mui/material/CssBaseline";
import Stack from "@mui/material/Stack";

import AnalyticsRoundedIcon from "@mui/icons-material/AnalyticsRounded";
import FileDownload from "@mui/icons-material/FileDownload";
import HelpRoundedIcon from "@mui/icons-material/HelpRounded";
import InfoRoundedIcon from "@mui/icons-material/InfoRounded";
import QuestionAnswer from "@mui/icons-material/QuestionAnswer";
import SettingsRoundedIcon from "@mui/icons-material/SettingsRounded";
import TerminalIcon from "@mui/icons-material/Terminal";

import type {} from "@mui/x-date-pickers/themeAugmentation";
import type {} from "@mui/x-charts/themeAugmentation";
import type {} from "@mui/x-data-grid/themeAugmentation";
import type {} from "@mui/x-tree-view/themeAugmentation";

import AppNavbar from "./widgets/AppNavbar";
import Header from "./widgets/Header";
import SideMenu from "./widgets/SideMenu";
import AppTheme from "./shared-theme/AppTheme";
import {
  chartsCustomizations,
  dataGridCustomizations,
  datePickersCustomizations,
  treeViewCustomizations,
} from "../theme/customizations";

const xThemeComponents = {
  ...chartsCustomizations,
  ...dataGridCustomizations,
  ...datePickersCustomizations,
  ...treeViewCustomizations,
};
const mainListItems = [
  { text: "Sources", icon: <FileDownload />, href: "sources" },
  { text: "Object queries", icon: <QuestionAnswer />, href: "queries" },
  { text: "Dashboard", icon: <AnalyticsRoundedIcon />, href: "dashboard" },
  { text: "Logs", icon: <TerminalIcon /> },
];

const secondaryListItems = [
  { text: "Settings", icon: <SettingsRoundedIcon /> },
  { text: "About", icon: <InfoRoundedIcon /> },
  { text: "Feedback", icon: <HelpRoundedIcon /> },
];

export default function LandingPage(props: { disableCustomTheme?: boolean }) {
  let navigate = useNavigate();
  const path = useLocation().pathname;
  let navIdx = 0;
  for (let i = 0; i < mainListItems.length; i++) {
    const href = mainListItems[i].href || "";
    if (href && path.indexOf(href) > -1) {
      navIdx = i;
      break;
    }
  }
  const [menuItemSelectedIdx, setMenuItemSelectedIdx] = useState(navIdx);

  const navigateToSection = (idx: number) => {
    setMenuItemSelectedIdx(idx);
    if (mainListItems[idx].href) {
      navigate(mainListItems[idx].href);
    }
  };

  return (
    <AppTheme {...props} themeComponents={xThemeComponents}>
      <CssBaseline enableColorScheme />
      <Box sx={{ display: "flex" }}>
        <SideMenu
          mainItems={mainListItems}
          secondaryItems={secondaryListItems}
          onSelected={(idx) => navigateToSection(idx)}
          menuItemSelectedIdx={menuItemSelectedIdx}
        />
        <AppNavbar pageTitle={mainListItems[menuItemSelectedIdx].text} />
        {/* Main content */}
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
            <Header />
            <Outlet />
          </Stack>
        </Box>
      </Box>
    </AppTheme>
  );
}
