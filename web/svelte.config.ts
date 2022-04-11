import type { Options } from "@sveltejs/vite-plugin-svelte";
import sveltePreprocess from "svelte-preprocess";
import { scssLegacyAliasImporter } from "./config/resolvers";

const production = process.env["NODE_ENV"] === "production";
const config: Options = {
  preprocess: sveltePreprocess({
    sourceMap: !production,
    scss: {
      importer: scssLegacyAliasImporter({
        "@": "./src",
        "~": "./node_modules",
      }),
      prependData: '@use "@styles/variables" as *;',
    },
  }),
};
export default config;
