import { defineConfig } from "vite";
import react from "@vitejs/plugin-react-swc";

// https://vite.dev/config/
export default defineConfig({
  plugins: [react()],
  server: {
    proxy: {
      "/rpc": {
        target: "http://localhost:8090",
        changeOrigin: true,
        secure: false,
      },
    },
  },
  build: {
    rollupOptions: {
      output: {
        manualChunks(id) {
          if (id.includes("node_modules")) {
            const libParts = id
              .toString()
              .split("node_modules/")[1]
              .toString()
              .split("/");
            if (libParts[0] === "@mui") {
              if (libParts[1].indexOf("x-") === 0) {
                return "@mui-" + libParts[1];
              }
              return "@mui";
            }
            return "lib";
          }
        },
      },
    },
  },
});
