import * as path from "path";
import * as fs from "fs/promises";
import type { Alias, ResolverFunction } from "vite";
import type { LegacyImporter, LegacyImporterResult } from "sass";
import type { PossiblyResolver } from ".";
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
  return function async(
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
