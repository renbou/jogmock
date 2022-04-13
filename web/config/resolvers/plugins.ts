import type { Alias, ResolverFunction } from "vite";
import type { LegacyImporter, LegacyImporterResult } from "sass";
import type { PossiblyResolver } from ".";
import * as path from "path";
import * as fs from "fs/promises";
import resolver from ".";

class banlistRegExp extends RegExp {
  banlist: RegExp[];

  constructor(pattern: string, banlist: RegExp[], flags?: string) {
    super(pattern, flags);
    this.banlist = banlist;
  }

  test(string: string): boolean {
    for (const ban of this.banlist) {
      if (ban.test(string)) {
        return false;
      }
    }
    return super.test(string);
  }
}

export type Aliases = {
  [alias: string]: string;
};

// aliasesResolver returns a mergingResolver which resolves
// using all of the given aliases
const aliasesResolver = (
  aliases: Aliases,
  ...resolvers: PossiblyResolver[]
): PossiblyResolver => {
  return resolver.merging(
    ...Object.getOwnPropertyNames(aliases).map((alias) =>
      resolver.alias(
        { alias, directory: path.resolve(aliases[alias]) },
        // resolve directories by hand first, then pass to vite's resolver
        ...resolvers
      )
    )
  );
};

// viteAliasResolver returns Vite Aliases initialized
// to resolve in a prettier and better way than by default.
// Use bypass to specify regexp of what not to alias
// (e.g. @vite or other problematic, dynamic modules)
export const viteAlias = (
  aliases: Aliases,
  exclude: RegExp[] = [/@vite/]
): Alias[] => {
  const r = aliasesResolver(
    aliases,
    resolver.pkgJson,
    resolver.jsIndex,
    resolver.identity
  );
  const viteAliasResolver: ResolverFunction = async function (
    importee,
    importer,
    resolveOptions
  ) {
    const resolve = (src: string) =>
      this.resolve(
        src,
        importer,
        Object.assign({ skipSelf: true }, resolveOptions)
      );

    const possibleResolutions = await Promise.resolve(r(importee));
    if (possibleResolutions !== undefined && possibleResolutions.length > 0) {
      for (const possibleResolution of possibleResolutions) {
        const resolved = await resolve(possibleResolution);
        if (Boolean(resolved)) {
          return resolved;
        }
      }
    }
    return (await resolve(importee)) || { id: importee };
  };

  return Object.getOwnPropertyNames(aliases).map((alias) => {
    // normal regexp match except for bypasses
    return {
      find: new banlistRegExp(`^${alias}(.+)`, exclude),
      replacement: `${alias}$1`,
      customResolver: viteAliasResolver,
    };
  });
};

export const scssLegacyAliasImporter = (
  aliases: Aliases
): LegacyImporter<"async"> => {
  const r = aliasesResolver(aliases, resolver.sass, resolver.pkgJson);
  return function (
    url: string,
    prev: string,
    done: (result: LegacyImporterResult) => void
  ) {
    Promise.resolve(r(url)).then((possibleResolutions) => {
      if (possibleResolutions !== undefined && possibleResolutions.length > 0) {
        done({ file: possibleResolutions[0] });
        return;
      }
      done(null);
    });
  };
};

// Type of data stored in a variables.json
type cssVariables = {
  [key: string]: string | string[] | cssVariables;
};

export function scssLegacyJsonImporter(
  url: string,
  prev: string,
  done: (result: LegacyImporterResult) => void
): void {
  // Custom scheme is needed because otherwise Vite breaks?? for some reason
  if (!url.startsWith("json:")) {
    done(null);
    return;
  }
  url = url.slice("json:".length);

  fs.readFile(path.join(path.dirname(prev), `${url}.json`)).then((data) => {
    const json = JSON.parse(data.toString());

    const camelToKebab = (s: string): string => {
      return s
        .replace(/([A-Z]+)/g, " $1")
        .split(" ")
        .filter((s) => s !== "")
        .map((s) => s.toLowerCase())
        .join("-");
    };

    let variables = "";
    // Convert all variables in json to kebab-case: value
    const generateVars = (prefix: string, map: cssVariables) => {
      for (const key in map) {
        const kebabKey = `${prefix}${prefix && "-"}${camelToKebab(key)}`;
        const value = map[key] as cssVariables;
        if (typeof value === "string") {
          variables += `$${kebabKey}: ${map[key]};\n`;
        } else if (value instanceof Array) {
          variables += `$${kebabKey}: ${value.join(", ")};\n`;
        } else {
          generateVars(kebabKey, value);
        }
      }
    };
    generateVars(camelToKebab(path.parse(url).name), json as cssVariables);
    done({ contents: variables });
  });
}
