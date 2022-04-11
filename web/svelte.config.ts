import type { Options } from "@sveltejs/vite-plugin-svelte";
import sveltePreprocess from "svelte-preprocess";
import {
  scssLegacyAliasImporter,
  scssLegacyJsonImporter,
} from "./config/resolvers";
import { duringDev } from "./config";

const config: Options = {
  preprocess: sveltePreprocess({
    sourceMap: duringDev,
    scss: {
      importer: [
        scssLegacyAliasImporter({
          "@": "./src",
          "~": "./node_modules",
        }),
        scssLegacyJsonImporter,
      ],
      prependData: '@use "@styles/variables" as *;',
    },
  }),
};
export default config;
