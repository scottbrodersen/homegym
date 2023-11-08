import { defineConfig } from 'vite';
import vue from '@vitejs/plugin-vue';
import { quasar, transformAssetUrls } from '@quasar/vite-plugin';

// https://vitejs.dev/config/
export default defineConfig({
  base: '/homegym/home/',
  mode: 'development',
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
    sourcemap: true,
    minify: false,
    //    watch: {},
    mode: 'development',
    emptyOutDir: true,
    rollupOptions: {
      output: {
        entryFileNames: `assets/entry-[name].js`,
        chunkFileNames: `assets/chunk-[name].js`,
        assetFileNames: `assets/asset-[name].[ext]`,
        compact: false,
        preserveModules: false,
      },
      preserveEntrySignatures: true,
    },
  },
});
