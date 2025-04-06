import { useState } from "react";

import "../../assets/styles.scss";

import { Outlet, useLocation, useNavigate } from "react-router-dom";

import { alpha } from "@mui/material/styles";
import Box from "@mui/material/Box";
import CssBaseline from "@mui/material/CssBaseline";
import Stack from "@mui/material/Stack";

import type {} from "@mui/x-date-pickers/themeAugmentation";
import type {} from "@mui/x-charts/themeAugmentation";
import type {} from "@mui/x-data-grid/themeAugmentation";
import type {} from "@mui/x-tree-view/themeAugmentation";

import { mainListItems, secondaryListItems } from "./rootmap";
import LogDrawer from "./logs/LogDrawer";
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

export default function HostPage(props: { disableCustomTheme?: boolean }) {
  let navigate = useNavigate();
  const path = useLocation().pathname;
  let navIdx = 1;
  for (let i = 0; i < mainListItems.length; i++) {
    const pth = mainListItems[i].path || "";
    if (pth && path.indexOf(pth) > -1) {
      navIdx = i;
      break;
    }
  }
  const [openLogPane, setOpenLogPane] = useState(
    localStorage.getItem("app-logpane") === "open"
  );
  const [menuItemSelectedIdx, setMenuItemSelectedIdx] = useState(navIdx);

  const navigateToSection = (idx: number) => {
    setMenuItemSelectedIdx(idx);
    if (mainListItems[idx].path) {
      navigate(mainListItems[idx].path);
    }
  };

  const secondaryItemClicked = (idx: number) => {
    if (idx === 0) {
      setOpenLogPanePersist(!openLogPane);
    }
  };

  const setOpenLogPanePersist = (b: boolean) => {
    localStorage.setItem("app-logpane", b ? "open" : "closed");
    setOpenLogPane(b);
  };

  return (
    <AppTheme {...props} themeComponents={xThemeComponents}>
      <CssBaseline enableColorScheme />
      <Box sx={{ display: "flex" }}>
        <SideMenu
          mainItems={mainListItems}
          secondaryItems={secondaryListItems}
          onSelected={(idx) => navigateToSection(idx)}
          onSecondarySelected={(idx) => secondaryItemClicked(idx)}
          menuItemSelectedIdx={menuItemSelectedIdx}
        />
        <AppNavbar
          mainItems={mainListItems}
          secondaryItems={secondaryListItems}
          pageTitle={mainListItems[menuItemSelectedIdx].text}
          onSelected={(idx) => navigateToSection(idx)}
          onSecondarySelected={(idx) => secondaryItemClicked(idx)}
        />
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
            spacing={4}
            sx={{
              alignItems: "center",
              mx: 3,
              pb: 5,
              mt: { xs: 8, md: 0 },
            }}
          >
            <Header pageTitle={mainListItems[menuItemSelectedIdx].text || ""} />
            <Outlet />
          </Stack>
        </Box>
      </Box>
      <LogDrawer open={openLogPane} setOpen={setOpenLogPanePersist} />
    </AppTheme>
  );
}
