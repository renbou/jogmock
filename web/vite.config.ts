import { defineConfig } from "vite";
import { svelte } from "@sveltejs/vite-plugin-svelte";
import { createHtmlPlugin as html } from "vite-plugin-html";
import { viteAlias } from "./config/resolvers";
import postcss from "./postcss.config";

export default defineConfig(({ mode }) => {
  return {
    plugins: [
      svelte(), // Svelte plugin options are contained within svelte.config.ts
      html({ minify: true }),
    ],
    resolve: {
      alias: viteAlias({
        "@": "./src",
        "~": "./node_modules",
      }),
    },
    css: {
      postcss, // PostCss options are contained within postcss.config.ts
      preprocessorOptions: {
        scss: {
          charset: false, // Remove useless CSS charsets
        },
      },
    },
  };
});
