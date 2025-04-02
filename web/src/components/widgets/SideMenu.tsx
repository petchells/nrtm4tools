import { styled } from "@mui/material/styles";
import MuiDrawer, { drawerClasses } from "@mui/material/Drawer";
import Box from "@mui/material/Box";
import Divider from "@mui/material/Divider";
import Typography from "@mui/material/Typography";

import BuildCircleIcon from "@mui/icons-material/BuildCircle";
import HardwareIcon from "@mui/icons-material/Hardware";
import PlumbingIcon from "@mui/icons-material/Plumbing";
import SquareFootIcon from "@mui/icons-material/SquareFoot";

import MenuContent from "./MenuContent";
import Stack from "@mui/material/Stack";

const drawerWidth = 180;

const Drawer = styled(MuiDrawer)({
  width: drawerWidth,
  flexShrink: 0,
  boxSizing: "border-box",
  mt: 10,
  [`& .${drawerClasses.paper}`]: {
    width: drawerWidth,
    boxSizing: "border-box",
  },
});

const BuildCircleIcon90 = styled(BuildCircleIcon)`
  transform: rotate(90deg);
`;
export default function SideMenu(props: {
  mainItems: any[];
  secondaryItems: any[];
  menuItemSelectedIdx: number;
  onSelected: (idx: number) => void;
  onSecondarySelected: (idx: number) => void;
}) {
  return (
    <Drawer
      variant="permanent"
      sx={{
        display: { xs: "none", md: "block" },
        [`& .${drawerClasses.paper}`]: {
          backgroundColor: "background.paper",
        },
      }}
    >
      <Box
        sx={{
          display: "flex",
          mt: "calc(var(--template-frame-height, 0px) + 4px)",
          p: 1.5,
          pl: 2.5,
        }}
      >
        <Stack>
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
      <MenuContent
        mainItems={props.mainItems}
        secondaryItems={props.secondaryItems}
        onSelected={(idx) => props.onSelected(idx)}
        onSecondarySelected={(idx) => props.onSecondarySelected(idx)}
        menuItemSelectedIdx={props.menuItemSelectedIdx}
      />
    </Drawer>
  );
}
