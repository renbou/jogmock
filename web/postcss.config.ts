import type * as PostCss from "postcss";
import purgecss from "@fullhuman/postcss-purgecss";
import csso from "postcss-csso";
import tailwindcss from "tailwindcss";
import tailwindConfig from "./tailwind.config";
import autoprefixer from "autoprefixer";
import { defaultExtractor } from "./config/purgecss";
import { duringProd } from "./config";

type PostcssConfig = PostCss.ProcessOptions & {
  plugins?: PostCss.Plugin[];
};
const config: PostcssConfig = {
  plugins: (() => {
    const plugins: PostCss.Plugin[] = [
      tailwindcss(tailwindConfig),
      autoprefixer(),
    ];
    if (duringProd) {
      plugins.push(
        purgecss({
          content: ["src/**/*.svelte"],
          defaultExtractor,
          // Keep html, body which are only in index.html as well as Svelte's scoped classes
          safelist: ["html", "body", /svelte-/],
        }) as PostCss.Plugin,
        csso({
          comments: false,
        })
      );
    }
    return plugins;
  })(),
};
export default config;
