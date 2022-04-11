import type { Options } from "@sveltejs/vite-plugin-svelte";
import sveltePreprocess from "svelte-preprocess";
import { scssLegacyAliasImporter } from "./config/resolvers";
import { isDev } from "./config";

const config: Options = {
  preprocess: sveltePreprocess({
    sourceMap: isDev(),
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
