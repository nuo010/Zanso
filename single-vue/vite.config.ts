import { defineConfig } from 'vite';
import vue from '@vitejs/plugin-vue';
import WindiCSS from 'vite-plugin-windicss';
import path from 'path';

const resolve = (dir: string) => path.resolve(__dirname, dir);

const alias: Record<string, string> = {
  '@': resolve('src'),
};
export default defineConfig({
  base: './',
  server: {
    port: 5173,
    // 监听所有网段
    host: true,
    // 允许所有host
    allowedHosts: true
  },
  plugins: [vue(), WindiCSS()],
  resolve: {
    alias,
  },
});
