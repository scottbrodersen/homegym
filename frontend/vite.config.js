/// <reference types="vitest" />
import { defineConfig } from 'vite';
import vue from '@vitejs/plugin-vue';
import { quasar } from '@quasar/vite-plugin';
import VueDevTools from 'vite-plugin-vue-devtools';

export default defineConfig(({ command, mode }) => {
  const config = {
    base: '/homegym/home/',
    plugins: [
      VueDevTools(),
      vue(),
      quasar({
        sassVariables:
          '/Users/scottbrodersen/Documents/code/homegym/frontend/src/quasar-variables.sass',
        extras: ['material-icons'],
      }),
    ],
    css: {
      modules: {
        localsConvention: 'dashesOnly',
      },
    },
    build: {
      outDir: '../server/secured/dist',
      manifest: 'homegym-manifest.json',
      emptyOutDir: true,
      input: 'main.js',
    },
    test: {
      reporters: ['verbose', 'json'],
      outputFile: './test_results/test-output.json',
      environment: 'jsdom',
      globals: true,
    },
  };

  if (mode != 'prod') {
    config.base = '/homegym/home';

    config.mode = 'development';
    config.build.mode = 'development';
    config.build.watch = {};
    config.build.minify = false;
    config.build.sourcemap = true;
    config.build.rollupOptions = {
      output: {
        compact: false,
        entryFileNames: `assets/entry-[name].js`,
        chunkFileNames: `assets/chunk-[name].js`,
        assetFileNames: `assets/asset-[name].[ext]`,
      },
    };
    config.server = {
      hmr: {
        host: '127.0.0.1',
        overlay: true,
      },
      proxy: {
        '/homegym/login': {
          target: 'http://127.0.0.1:3000',
          changeOrigin: true,
          secure: false,
          ws: true,
          configure: defaultProxyConfig,
        },
        '/homegym/login/': {
          target: 'http://127.0.0.1:3000',
          changeOrigin: true,
          secure: false,
          ws: true,
          configure: defaultProxyConfig,
        },
        '/homegym/api': {
          target: 'http://127.0.0.1:3000',
          changeOrigin: true,
          secure: false,
          ws: true,
          configure: defaultProxyConfig,
        },
      },
    };
  } else if (mode == 'prod') {
    config.mode = 'production';
    config.build.minify = 'esbuild';
    config.optimizeDeps = { force: true };
  }

  return config;
});

const defaultProxyConfig = (proxy, _options) => {
  proxy.on('error', (err, _req, _res) => {
    console.log('proxy error', err);
  });
  proxy.on('proxyReq', (proxyReq, req, _res) => {
    console.log(req.method, req.url);
  });
  proxy.on('proxyRes', (proxyRes, req, _res) => {
    console.log(proxyRes.statusCode, req.url);
  });
};
