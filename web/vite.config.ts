import type * as Postcss from "postcss";
import { defineConfig } from "vite";
import { svelte } from "@sveltejs/vite-plugin-svelte";
import { createHtmlPlugin as html } from "vite-plugin-html";
import { viteAlias } from "./config/resolvers";
import purgecss from "@fullhuman/postcss-purgecss";
import purgeSvelte from "./config/purgecss";
import cssnano from "cssnano";

const postCssPlugins = (prod: boolean): Postcss.Plugin[] => {
  return !prod
    ? []
    : [
        purgecss({
          content: ["src/**/*.svelte"],
          extractors: [{ extensions: ["svelte"], extractor: purgeSvelte }],
          // Keep html, body which are only in index.html as well as Svelte's scoped classes
          safelist: ["html", "body", /svelte-/],
        }) as Postcss.Plugin,
        cssnano({
          preset: ["default", { discardComments: { removeAll: true } }],
        }),
      ];
};

export default defineConfig(({ mode }) => {
  const prod = mode === "production";
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
      postcss: {
        plugins: postCssPlugins(prod),
      },
      preprocessorOptions: {
        scss: {
          charset: false, // remove useless CSS charsets
        },
      },
    },
  };
});
