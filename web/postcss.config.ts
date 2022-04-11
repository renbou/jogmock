import type * as Postcss from "postcss";
import purgecss from "@fullhuman/postcss-purgecss";
import { defaultExtractor } from "./config/purgecss";
import cssnano from "cssnano";
import { isDev } from "./config";

type PostcssConfig = Postcss.ProcessOptions & {
  plugins?: Postcss.Plugin[];
};

const config: PostcssConfig = {
  plugins: isDev()
    ? []
    : [
        purgecss({
          content: ["src/**/*.svelte"],
          defaultExtractor,
          // Keep html, body which are only in index.html as well as Svelte's scoped classes
          safelist: ["html", "body", /svelte-/],
        }) as Postcss.Plugin,
        cssnano({
          preset: ["default", { discardComments: { removeAll: true } }],
        }),
      ],
};
export default config;
