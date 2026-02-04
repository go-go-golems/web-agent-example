import path from 'path';
import react from '@vitejs/plugin-react';
import { defineConfig } from 'vite';

const webchatRoot = path.resolve(__dirname, '../../pinocchio/cmd/web-chat/web/src');

export default defineConfig({
  plugins: [react()],
  resolve: {
    alias: {
      '@pwchat': webchatRoot,
    },
  },
  base: './',
  build: {
    outDir: '../cmd/web-agent-example/static/dist',
    emptyOutDir: true,
    sourcemap: true,
  },
  server: {
    port: 5174,
    proxy: {
      '/chat': { target: process.env.VITE_BACKEND_ORIGIN ?? 'http://localhost:8080', changeOrigin: true },
      '/api': { target: process.env.VITE_BACKEND_ORIGIN ?? 'http://localhost:8080', changeOrigin: true },
      '/ws': { target: process.env.VITE_BACKEND_ORIGIN ?? 'http://localhost:8080', ws: true, changeOrigin: true },
      '/hydrate': { target: process.env.VITE_BACKEND_ORIGIN ?? 'http://localhost:8080', changeOrigin: true },
      '/timeline': { target: process.env.VITE_BACKEND_ORIGIN ?? 'http://localhost:8080', changeOrigin: true },
    },
  },
});
