import type tailwindcss from "tailwindcss";
import styleConfig from "../src/styles/variables.json";

export type TailwindConfig = Exclude<Parameters<typeof tailwindcss>[0], string>;
const config: TailwindConfig = {
  content: ["src/**/*.svelte"],
  theme: {
    extend: styleConfig,
  },
  plugins: [],
};
export default config;
