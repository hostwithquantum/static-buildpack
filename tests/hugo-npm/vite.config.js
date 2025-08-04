import { defineConfig } from 'vite'
import tailwindcss from '@tailwindcss/vite'

export default defineConfig({
    plugins: [
        tailwindcss({
            content: ['./content/**/*.{html,md}', './layouts/**/*.html']
        }),
    ],
    build: {
        outDir: 'static',
        rollupOptions: {
            input: 'assets/css/main.css',
            output: {
                assetFileNames: 'style.css',
                entryFileNames: '[name].js',
                chunkFileNames: '[name].js'
            }
        },
    }
})