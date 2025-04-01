import { useState } from "react";

import AppBar from "@mui/material/AppBar";
import Box from "@mui/material/Box";
import Toolbar from "@mui/material/Toolbar";
import IconButton from "@mui/material/IconButton";
import Typography from "@mui/material/Typography";
import Menu from "@mui/material/Menu";
import MenuIcon from "@mui/icons-material/Menu";
import Container from "@mui/material/Container";
import Button from "@mui/material/Button";
import Tooltip from "@mui/material/Tooltip";
import MenuItem from "@mui/material/MenuItem";

import CloseIcon from "@mui/icons-material/Close";
import TerminalIcon from "@mui/icons-material/Terminal";

const logLevel = ["Error", "Warning", "Info", "Debug"];

interface FrameToolbarProps {
  status: string;
  setOpen: (b: boolean) => void;
}

export default function FrameToolbar({ status, setOpen }: FrameToolbarProps) {
  const [anchorElUser, setAnchorElUser] = useState<null | HTMLElement>(null);

  const handleOpenUserMenu = (event: React.MouseEvent<HTMLElement>) => {
    setAnchorElUser(event.currentTarget);
  };

  const levelClickHandlerWrapper = (lvl: string) => () => {
    console.log("clicked", lvl);
  };

  const handleCloseUserMenu = () => {
    setAnchorElUser(null);
  };

  const handleClosePanel = () => {
    setAnchorElUser(null);
    setOpen(false);
  };

  return (
    <AppBar position="sticky">
      <Container maxWidth="xl">
        <Toolbar disableGutters>
          <Box sx={{ flexGrow: 0, display: { xs: "flex", sm: "none" }, mr: 1 }}>
            <Tooltip title="Open settings">
              <IconButton onClick={handleOpenUserMenu} sx={{ p: 0 }}>
                <MenuIcon />
              </IconButton>
            </Tooltip>
            <Menu
              sx={{ mt: "45px" }}
              id="menu-appbar"
              anchorEl={anchorElUser}
              anchorOrigin={{
                vertical: "top",
                horizontal: "right",
              }}
              keepMounted
              transformOrigin={{
                vertical: "top",
                horizontal: "right",
              }}
              open={Boolean(anchorElUser)}
              onClose={handleCloseUserMenu}
            >
              {logLevel.map((lvl) => (
                <MenuItem key={lvl} onClick={handleCloseUserMenu}>
                  <Typography sx={{ textAlign: "center" }}>{lvl}</Typography>
                </MenuItem>
              ))}
            </Menu>
          </Box>
          <TerminalIcon sx={{ display: { xs: "none", sm: "flex" }, mr: 1 }} />
          <Typography
            variant="h6"
            noWrap
            sx={{
              mr: 2,
              fontWeight: 700,
            }}
          >
            Logs
          </Typography>
          <Box sx={{ flexGrow: 1, display: "flex" }}>
            <Box sx={{ display: { xs: "none", sm: "flex" } }}>
              {logLevel.map((lvl) => (
                <Button
                  key={lvl}
                  onClick={levelClickHandlerWrapper(lvl)}
                  sx={{ my: 0, color: "white", display: "block" }}
                >
                  {lvl}
                </Button>
              ))}
            </Box>
          </Box>
          <Typography
            variant="body1"
            noWrap
            sx={{
              mr: 2,
              color: status === "Open" ? "white" : "red",
            }}
          >
            WS
          </Typography>
          <Box sx={{ flexGrow: 0, alignItems: "right" }}>
            <Tooltip title="Close panel">
              <IconButton onClick={handleClosePanel} sx={{ p: 0 }}>
                <CloseIcon />
              </IconButton>
            </Tooltip>
          </Box>
        </Toolbar>
      </Container>
    </AppBar>
  );
}
