import { styled } from "@mui/material/styles";
import Button from "@mui/material/Button";
import Box from "@mui/material/Box";
import Divider from "@mui/material/Divider";
import Drawer, { drawerClasses } from "@mui/material/Drawer";
import Stack from "@mui/material/Stack";
import Typography from "@mui/material/Typography";

import BuildCircleIcon from "@mui/icons-material/BuildCircle";
import HardwareIcon from "@mui/icons-material/Hardware";
import LogoutRoundedIcon from "@mui/icons-material/LogoutRounded";
import PlumbingIcon from "@mui/icons-material/Plumbing";
import SquareFootIcon from "@mui/icons-material/SquareFoot";

import MenuContent from "./MenuContent";
import { MenuItem } from "./widgettypes";

const BuildCircleIcon90 = styled(BuildCircleIcon)`
  transform: rotate(90deg);
`;

interface SideMenuMobileProps {
  open: boolean | undefined;
  toggleDrawer: (newOpen: boolean) => () => void;
  onSelected: (idx: number) => void;
  onSecondarySelected: (idx: number) => void;
  mainItems: MenuItem[];
  secondaryItems: MenuItem[];
}

export default function SideMenuMobile({
  open,
  toggleDrawer,
  onSelected,
  onSecondarySelected,
  mainItems,
  secondaryItems,
}: SideMenuMobileProps) {
  return (
    <Drawer
      anchor="right"
      open={open}
      onClose={toggleDrawer(false)}
      sx={{
        zIndex: (theme) => theme.zIndex.drawer + 1,
        [`& .${drawerClasses.paper}`]: {
          backgroundImage: "none",
          backgroundColor: "background.paper",
        },
      }}
    >
      <Stack
        sx={{
          maxWidth: "70dvw",
          height: "100%",
        }}
      >
        <Box
          sx={{
            display: "flex",
            mt: "calc(var(--template-frame-height, 0px) + 4px)",
            p: 1.5,
          }}
        >
          <Stack alignItems="center">
            <Typography
              component="h2"
              variant="h3"
              sx={{ color: "text.secondary" }}
            >
              NRTM4
            </Typography>
            <Typography
              component="h2"
              variant="h3"
              sx={{ color: "text.secondary" }}
            >
              <HardwareIcon />
              <BuildCircleIcon />
              <BuildCircleIcon90 />
              <SquareFootIcon />
              <PlumbingIcon />
            </Typography>
          </Stack>
        </Box>
        <Divider />
        <Stack sx={{ flexGrow: 1 }}>
          <MenuContent
            mainItems={mainItems}
            secondaryItems={secondaryItems}
            menuItemSelectedIdx={0}
            onSelected={(idx) => onSelected(idx)}
            onSecondarySelected={(idx) => onSecondarySelected(idx)}
          />
          <Divider />
        </Stack>
        <Stack sx={{ p: 2 }}>
          <Button
            variant="outlined"
            fullWidth
            startIcon={<LogoutRoundedIcon />}
          ></Button>
        </Stack>
      </Stack>
    </Drawer>
  );
}
