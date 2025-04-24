// import { defineConfig } from 'vite'
// import { svelte } from '@sveltejs/vite-plugin-svelte'
//
// // https://vite.dev/config/
// export default defineConfig({
//   plugins: [svelte()],
// })

import { defineConfig } from 'vite';
import { svelte } from '@sveltejs/vite-plugin-svelte';
import { resolve } from 'path';

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [svelte()],

  // Set the base URL (important for path resolution when embedded in Go)
  base: '/',

  // Configure the build output to match Go's embed requirements
  build: {
    // Output to the static directory at project root
    outDir: resolve(__dirname, '../static'),

    // Make sure paths are relative, not absolute
    assetsDir: 'assets',

    // For better debugging in development
    sourcemap: true,

    // Clean the output directory before each build
    emptyOutDir: true
  },

  // During development, proxy API requests to the Go backend
  server: {
    cors: true,
  },

  // Resolve TypeScript path aliases
  resolve: {
    alias: {
      '$lib': resolve(__dirname, './src/lib')
    }
  }
});