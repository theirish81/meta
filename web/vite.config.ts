import tailwindcss from '@tailwindcss/vite';
import { defineConfig } from 'vite';
import { svelte } from '@sveltejs/vite-plugin-svelte';

// https://vite.dev/config/
export default defineConfig({
    plugins: [tailwindcss(), svelte()],
    base: "/web",
    server: {
        fs: {
            strict: false
        },
        proxy: {
            // Proxy API calls starting with /api to your backend
            '/api': {
                target: 'http://localhost:8080', // replace with your backend URL
                changeOrigin: false,
            },
        },
    }
});
