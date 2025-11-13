/// <reference types="vitest" />
import { defineConfig } from 'vite';
import react from '@vitejs/plugin-react';
import tailwindcss from '@tailwindcss/vite';
import { visualizer } from 'rollup-plugin-visualizer';



// https://vite.dev/config/
export default defineConfig(({ mode }) => ({
    plugins: [
        react(),
        tailwindcss(),
        // Visualizer plugin to analyze bundle size
        mode === 'analyze' && visualizer({
        open: true,
        filename: 'dist/bundle-stats.html',
        gzipSize: true,
        brotliSize: true,
        }),
    ],
    build: {
        rollupOptions: {
            output: {
                manualChunks(id) {
                    if (id.includes('node_modules')) {
                        if (id.includes('/react/') || id.includes('/react-dom/') || id.includes('/react-router/') || id.includes('/react-router-dom/')) {
                            // React in separate chunks
                            return 'react-vendor';
                        } else if (id.includes('/@tanstack/')) {
                            // TanStack libraries in separate chunks
                            return 'tanstack-vendor';
                        }
                    }
                    return 'vendor';
                }
            }
        },
        chunkSizeWarningLimit: 1000, // Increase warning limit for map chunks
    },
    test: {
        globals: true,
        environment: 'jsdom',
        setupFiles: './src/test/setup.ts',
    },
}));
