import { defineConfig } from 'vite';
import react from '@vitejs/plugin-react';

// Define the development port. You can change it to match your Vue.js configuration.
const devPort = 9528;

// Define the development proxy. You can change the target to match your Vue.js proxy configuration.
const devProxy = {
  '/api': {
    target: 'http://127.0.0.1:8888', // Used during local development to connect to the backend interface
    changeOrigin: true,
  },
};

export default defineConfig({
  plugins: [react()],

  // Set the base URL for your deployment, if needed.
  // base: '/',

  server: {
    // Configure the development server to match the Vue.js configuration.
    port: devPort,
    open: true,
    proxy: devProxy,
  },

  build: {
    // Output directory.
    outDir: 'dist',
  },
});