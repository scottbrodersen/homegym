import { defineConfig } from 'vite';
import vue from '@vitejs/plugin-vue';
import { quasar, transformAssetUrls } from '@quasar/vite-plugin';

export default defineConfig(({ command, mode }) => {
  const config = {
    base: '/homegym/home/',
    plugins: [
      vue({ template: { transformAssetUrls } }),
      quasar({
        sassVariables: 'src/quasar-variables.sass',
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
    },
  };

  if (mode == 'dev') {
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
  } else if (mode == 'prod') {
    config.mode = 'production';
    config.build.minify = 'esbuild';
    config.optimizeDeps = { force: true };
  }

  return config;
});
