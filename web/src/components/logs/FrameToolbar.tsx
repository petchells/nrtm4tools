import { useState } from "react";

import AppBar from "@mui/material/AppBar";
import Box from "@mui/material/Box";
import Button from "@mui/material/Button";
import ButtonGroup from "@mui/material/ButtonGroup";
import Container from "@mui/material/Container";
import Toolbar from "@mui/material/Toolbar";
import IconButton from "@mui/material/IconButton";
import Menu from "@mui/material/Menu";
import MenuItem from "@mui/material/MenuItem";
import Tooltip from "@mui/material/Tooltip";
import Typography from "@mui/material/Typography";

import CircleIcon from "@mui/icons-material/Circle";
import CloseIcon from "@mui/icons-material/Close";
import LeakAddIcon from "@mui/icons-material/LeakAdd";
import LeakRemoveIcon from "@mui/icons-material/LeakRemove";
import MenuIcon from "@mui/icons-material/Menu";
import TerminalIcon from "@mui/icons-material/Terminal";
import WatchIcon from "@mui/icons-material/Watch";
import WatchOffIcon from "@mui/icons-material/WatchOff";

import { ToolbarCommand } from "./model";
import Stack from "@mui/material/Stack";

const logLevel = ["Error", "Info", "Verbose", "Fine"];

interface FrameToolbarProps {
  status: string;
  scrollBottom: boolean;
  toolbarClick: (cmd: ToolbarCommand, ...args: any) => void;
}

export default function FrameToolbar({
  status,
  scrollBottom,
  toolbarClick,
}: FrameToolbarProps) {
  const [anchorElUser, setAnchorElUser] = useState<null | HTMLElement>(null);

  const handleOpenUserMenu = (event: React.MouseEvent<HTMLElement>) => {
    setAnchorElUser(event.currentTarget);
  };

  const levelClickHandlerWrapper = (lvl: number) => () => {
    setAnchorElUser(null);
    toolbarClick(ToolbarCommand.logLevel, lvl);
  };

  const handleCloseUserMenu = () => {
    setAnchorElUser(null);
  };

  const handleClosePanel = () => {
    setAnchorElUser(null);
    toolbarClick(ToolbarCommand.closeLogPane);
  };

  const colours = ["#cc0000", "#6633ff", "#47c18e", "#78bcff"];

  const levelButton = (lvl: string) => {
    const c = colours[logLevel.indexOf(lvl)];
    return (
      <Tooltip title={lvl} key={lvl}>
        <Button
          sx={{ color: c }}
          onClick={levelClickHandlerWrapper(logLevel.indexOf(lvl))}
        >
          <CircleIcon />
        </Button>
      </Tooltip>
    );
  };

  return (
    <AppBar
      position="sticky"
      sx={{ bgcolor: "background.paper", color: "text.primary" }}
    >
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
                <MenuItem
                  key={lvl}
                  onClick={levelClickHandlerWrapper(logLevel.indexOf(lvl))}
                >
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
              <ButtonGroup size="small" aria-label="Small button group">
                {logLevel.map((lvl) => levelButton(lvl))}
              </ButtonGroup>
            </Box>
          </Box>
          <Stack direction="row" spacing={1}>
            <Tooltip title="Watch the tail of the log">
              <IconButton
                size="small"
                onClick={() =>
                  toolbarClick(ToolbarCommand.scrollBottom, !scrollBottom)
                }
              >
                {scrollBottom ? <WatchIcon /> : <WatchOffIcon />}
              </IconButton>
            </Tooltip>
            <Tooltip title="Web socket status. Click to reconnect">
              <IconButton
                size="small"
                onClick={() => toolbarClick(ToolbarCommand.reconnectWS)}
              >
                {status === "Open" ? (
                  <LeakAddIcon />
                ) : (
                  <LeakRemoveIcon color="error" sx={{ cursor: "pointer" }} />
                )}
              </IconButton>
            </Tooltip>
            <Box sx={{ flexGrow: 0, alignItems: "right", ml: 1 }}>
              <Tooltip title="Close panel">
                <IconButton size="small" onClick={handleClosePanel}>
                  <CloseIcon />
                </IconButton>
              </Tooltip>
            </Box>
          </Stack>
        </Toolbar>
      </Container>
    </AppBar>
  );
}
