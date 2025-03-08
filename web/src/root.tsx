import { createBrowserRouter } from "react-router-dom";
import ErrorPage from "./error-page";
import HostPage from "./components/HostPage";
import { mainListItems } from "./components/rootmap";

const router = createBrowserRouter([
  {
    path: "/",
    element: <HostPage />,
    errorElement: <ErrorPage />,

    children: mainListItems,
  },
]);

export default router;
