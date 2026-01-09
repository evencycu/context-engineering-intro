import path from 'path';
import { defineConfig, loadEnv } from 'vite';
import react from '@vitejs/plugin-react';

export default defineConfig(({ mode }) => {
    const env = loadEnv(mode, '.', '');
    
    // Backend API server URL (default: localhost:8080 for local development)
    const BACKEND_URL = env.VITE_BACKEND_URL || 'http://localhost:8080';
    
    return {
      server: {
        port: 3000,
        host: '0.0.0.0',
        // Proxy API requests to backend server
        proxy: {
          // External API (public)
          '/api': {
            target: BACKEND_URL,
            changeOrigin: true,
            secure: false,
          },
          // Internal API (admin features)
          '/internal': {
            target: BACKEND_URL,
            changeOrigin: true,
            secure: false,
          },
          // Health check endpoint
          '/health': {
            target: BACKEND_URL,
            changeOrigin: true,
            secure: false,
          },
        },
      },
      plugins: [react()],
      define: {
        'process.env.API_KEY': JSON.stringify(env.GEMINI_API_KEY),
        'process.env.GEMINI_API_KEY': JSON.stringify(env.GEMINI_API_KEY),
        // Expose backend URL for runtime usage if needed
        'process.env.VITE_BACKEND_URL': JSON.stringify(BACKEND_URL),
      },
      resolve: {
        alias: {
          '@': path.resolve(__dirname, '.'),
        }
      }
    };
});
