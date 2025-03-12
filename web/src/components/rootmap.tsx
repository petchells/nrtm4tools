import { Navigate } from "react-router-dom";

import AnalyticsRoundedIcon from "@mui/icons-material/AnalyticsRounded";
import FileDownload from "@mui/icons-material/FileDownload";
import HelpRoundedIcon from "@mui/icons-material/HelpRounded";
import InfoRoundedIcon from "@mui/icons-material/InfoRounded";
import QuestionAnswer from "@mui/icons-material/QuestionAnswer";
import SettingsRoundedIcon from "@mui/icons-material/SettingsRounded";
import TerminalIcon from "@mui/icons-material/Terminal";

import Logs from "./logs/Logs";
import Sources from "./sources/Sources";
import MainGrid from "./dashboard/MainGrid";
import Queries from "./queries/Queries";
import { MenuItem } from "./widgets/widgettypes";

interface MenuRouter extends MenuItem {
  index?: boolean;
  element: any;
  path?: string;
}

export const mainListItems: MenuRouter[] = [
  {
    index: true,
    element: <Navigate to="sources" replace />,
  },
  {
    text: "Sources",
    icon: <FileDownload />,
    path: "sources",
    element: <Sources />,
  },
  {
    text: "Dashboard",
    icon: <AnalyticsRoundedIcon />,
    path: "dashboard",
    element: <MainGrid />,
  },
  { text: "Logs", icon: <TerminalIcon />, path: "logs", element: <Logs /> },
  {
    text: "Object queries",
    icon: <QuestionAnswer />,
    path: "queries",
    element: <Queries />,
  },
];


export const secondaryListItems: MenuItem[] = [
  { text: "Settings", icon: <SettingsRoundedIcon /> },
  { text: "About", icon: <InfoRoundedIcon /> },
  { text: "Feedback", icon: <HelpRoundedIcon /> },
];
