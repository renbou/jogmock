import type tailwindcss from "tailwindcss";

export type TailwindConfig = Exclude<Parameters<typeof tailwindcss>[0], string>;
const config: TailwindConfig = {
  content: ["src/**/*.svelte"],
  theme: {
    extend: {},
  },
  plugins: [],
};
export default config;
