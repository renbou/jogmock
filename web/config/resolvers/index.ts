import * as path from "path";
import * as fsp from "fs/promises";
import type { Stats } from "fs";

export type PossibleResolution = string[] | undefined;

// PossiblyResolver should return possible filepaths to which the source importee
// might resolve to. The source is considered resolved if at least one of the returned
// paths exists and is accessible.
export type PossiblyResolver = (
  source: string
) => Promise<PossibleResolution> | PossibleResolution;

// identity is a PossiblyResolver returning the
// passed source as a possible resolution.
function identity(source: string) {
  return [source];
}

// existing acts like identity but checks that
// the source is an existing path (file or directory).
// Returns maximum one entry containing the source argument.
async function existing(source: string, verifier?: (stats: Stats) => boolean) {
  return fsp
    .stat(source)
    .then((stats) => {
      return (verifier && verifier(stats)) || !verifier ? [source] : undefined;
    })
    .catch(() => false);
}
// existingFile is an existing resolver which additionally checks that the source is a file
async function existingFile(source: string) {
  return existing(source, (s) => s.isFile());
}
// existingFileResolver is an existing resolver which additionally checks that the source is a directory
async function existingDir(source: string) {
  return existing(source, (s) => s.isDirectory());
}

// extension returns a PossiblyResolver which tries to resolve files with the given extensions.
// It only resolves existing files, meaning it looks up if a file with the extension actually exists.
function extension(...extensions: string[]): PossiblyResolver {
  return async function (source: string) {
    const resolutions: string[] = [];
    for (const file of extensions.map((ext) => `${source}.${ext}`)) {
      if (await existingFile(file)) {
        resolutions.push(file);
      }
    }
    return resolutions.length > 0 ? resolutions : undefined;
  };
}

// jsExtension is an extension resolver for js files
const jsExtension: PossiblyResolver = extension("js", "cjs", "mjs");

// sassExtension is an extension resolver for sass files
const sassExtension: PossiblyResolver = extension("sass", "scss", "css");

// sass is a PossiblyResolver which tries to resolve
// sass files using the sass import rules.
async function sass(source: string) {
  if (/\.(css|scss|sass)$/.test(source)) {
    // we don't need to resolve anything if the source is already a sass file
    return [source];
  }

  const paths = [source, `${source}/index`, `${source}/_index`];
  if (!path.basename(source).startsWith("_")) {
    paths.push(`${path.dirname(source)}/_${path.basename(source)}`);
  }

  const resolutions = [];
  for (const path of paths) {
    const files = await sassExtension(path);
    if (files) {
      resolutions.push(...files);
    }
  }
  return resolutions;
}

// pkgjson is a PossiblyResolver which tries to resolve the source as a
// package with a package.json which contains the "main" field
async function pkgJson(source: string) {
  const packagePath = path.join(source, "package.json");
  try {
    if (await existingFile(packagePath)) {
      const pkgdata = JSON.parse((await fsp.readFile(packagePath)).toString());
      const pathToMain = pkgdata["main"]
        ? path.join(source, pkgdata["main"])
        : undefined;
      if (pathToMain) {
        return [pathToMain];
      }
    }
  } catch (err) {
    console.error(
      `Error resolving module package.json for path ${source}: ${
        err instanceof Error ? err.message : err
      }`
    );
  }
  return undefined;
}

// index is a PossiblyResolver which tries to resolve index.js, index.cjs and index.mjs
// files if source is a directory
async function jsIndex(source: string) {
  const indexResolutions = await jsExtension(path.join(source, "index"));
  if (indexResolutions) {
    return indexResolutions;
  }
  return undefined;
}

// merging returns a resolver which merges the results
// of resolutions through other resolvers
function merging(...resolvers: PossiblyResolver[]): PossiblyResolver {
  return async (source: string) => {
    const results: string[] = [];
    for (const resolver of resolvers) {
      await Promise.resolve(resolver(source)).then((resolution) => {
        if (resolution !== undefined) {
          results.push(...resolution);
        }
      });
    }
    return results;
  };
}

// alias creates a new resolver which resolves paths
// by unaliasing and passing them to other resolvers
function alias(
  {
    alias,
    directory,
  }: {
    alias: string;
    directory: string;
  },
  ...resolvers: PossiblyResolver[]
): PossiblyResolver {
  const merged = merging(...resolvers);
  return async (source: string): Promise<PossibleResolution> => {
    if (source.startsWith(alias) && source.length > alias.length) {
      return merged(path.join(directory, source.slice(alias.length)));
    }
    return undefined;
  };
}

export default {
  identity,
  existing,
  existingFile,
  existingDir,
  extension,
  jsExtension,
  sassExtension,
  merging,
  sass,
  pkgJson,
  jsIndex,
  alias,
};
export * from "./plugins";
