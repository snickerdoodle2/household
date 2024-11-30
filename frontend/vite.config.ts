import { sveltekit } from '@sveltejs/kit/vite';
import { defineConfig } from 'vite';

const serverPort = import.meta.env?.SERVER_PORT ?? '8080';
const serverHost = import.meta.env?.SERVER_HOST ?? 'localhost';

export default defineConfig({
    plugins: [sveltekit()],
    server: {
        proxy: {
            '/api': {
                target: `http://${serverHost}:${serverPort}`,
                ws: true,
            },
        },
    },
});
