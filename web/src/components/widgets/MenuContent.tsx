import List from "@mui/material/List";
import ListItem from "@mui/material/ListItem";
import ListItemButton from "@mui/material/ListItemButton";
import ListItemIcon from "@mui/material/ListItemIcon";
import ListItemText from "@mui/material/ListItemText";
import Stack from "@mui/material/Stack";
import { MenuItem } from "./widgettypes";

interface MenuContentProps {
  mainItems: MenuItem[];
  secondaryItems: any[];
  menuItemSelectedIdx: number;
  onSelected: (idx: number) => void;
}
export default function MenuContent({
  mainItems,
  secondaryItems,
  menuItemSelectedIdx,
  onSelected,
}: MenuContentProps) {
  // handle click on menu item
  const menuClicked = (idx: number) => () => onSelected(idx);

  return (
    <Stack sx={{ flexGrow: 1, p: 1, justifyContent: "space-between" }}>
      <List dense>
        {mainItems.map(
          (item, index) =>
            item.text && (
              <ListItem key={index} disablePadding sx={{ display: "block" }}>
                <ListItemButton
                  selected={index === menuItemSelectedIdx}
                  onClick={menuClicked(index)}
                >
                  <ListItemIcon>{item.icon}</ListItemIcon>
                  <ListItemText primary={item.text} />
                </ListItemButton>
              </ListItem>
            )
        )}
      </List>
      <List dense>
        {secondaryItems.map((item, index) => (
          <ListItem key={index} disablePadding sx={{ display: "block" }}>
            <ListItemButton>
              <ListItemIcon>{item.icon}</ListItemIcon>
              <ListItemText primary={item.text} />
            </ListItemButton>
          </ListItem>
        ))}
      </List>
    </Stack>
  );
}
