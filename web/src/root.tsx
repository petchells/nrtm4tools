import { createBrowserRouter, Navigate } from "react-router-dom";
import ErrorPage from "./error-page";
import HostPage from "./components/HostPage";
import Logs from "./components/logs/Logs";
import Sources from "./components/sources/Sources";
import MainGrid from "./components/dashboard/MainGrid";
import Queries from "./components/queries/Queries";

const router = createBrowserRouter([
  {
    path: "/",
    element: <HostPage />,
    errorElement: <ErrorPage />,

    children: [
      {
        index: true,
        element: <Navigate to="sources" replace />,
      },
      {
        path: "sources",
        element: <Sources />,
      },
      {
        path: "queries",
        element: <Queries />,
      },
      {
        path: "dashboard",
        element: <MainGrid />,
      },
      {
        path: "logs",
        element: <Logs />,
      },
    ],
  },
]);

export default router;
