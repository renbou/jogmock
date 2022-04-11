# Svelte template
Decent template for Svelte apps built with Vite and styled using Bulma.

To create a new project using this template `degit` can be used:
```bash
degit renbou/svelte-template svelte-app
cd svelte-app
```

## Features
- `Bulma` + `Sass` for styling with proper variable exports (all libraries' variables available)
- `Vite` for fast dev server and bundling with `Rollup`
- JS minification using Vite's default `esbuild`
- CSS minification using `PurgeCSS` and `cssnano`
- HTML minification using `html-minifier-terser`
- Alias resolvers providing allow `@`- and `~`-style imports for `src/` and `node_modules/` respectively (`~bulma`, `@components`) 

## How to
1. Install deps using `pnpm install` and run `pnpm start` to launch the `vite` dev server and `svelte-check` watcher.  

2. Compile and serve in production using `pnpm build` and `pnpm serve`. By default `serve` launches `sirv` in SPA mode, if the built app isn't an SPA then remove `--single` argument from the script.

## Optimizing built bundles
1. Rollup supports Code splitting through dynamic imports.
2. Make sure to use `<!-- prettier-ignore -->` to format html tags such as `p` for which Svelte preserves whitespace even when it isn't needed.