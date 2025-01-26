import { createBrowserRouter, Navigate } from "react-router-dom";
import ErrorPage from "./error-page";
import LandingPage from "./components/LandingPage";
import Sources from "./components/sources/Sources";
import MainGrid from "./components/dashboard/MainGrid";

const router = createBrowserRouter([
  {
    path: "/",
    element: <LandingPage />,
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
        element: <Sources />,
      },
      {
        path: "dashboard",
        element: <MainGrid />,
      },
    ],
  },
]);

export default router;
