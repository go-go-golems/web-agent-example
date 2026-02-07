import { defineConfig } from 'vite';
import react from '@vitejs/plugin-react';

export default defineConfig({
  plugins: [react()],
  base: './',
  build: {
    outDir: '../static/dist',
    emptyOutDir: true,
    sourcemap: true,
  },
  server: {
    port: 5174,
    proxy: {
      '/debug': 'http://localhost:8080',
      '/ws': { target: 'ws://localhost:8080', ws: true },
      '/chat': 'http://localhost:8080',
      '/timeline': 'http://localhost:8080',
      '/turns': 'http://localhost:8080',
      '/hydrate': 'http://localhost:8080',
      '/api': 'http://localhost:8080',
    },
  },
});
