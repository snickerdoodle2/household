import { sveltekit } from "@sveltejs/kit/vite";
import { defineConfig } from "vite";

const serverPort = import.meta.env?.SERVER_PORT ?? "8080";

export default defineConfig({
    plugins: [sveltekit()],
    server: {
        proxy: {
            "/api": {
                target: `http://localhost:${serverPort}`,
                ws: true,
            },
        },
    },
});
