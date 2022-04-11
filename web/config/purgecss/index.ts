import type purgecss from "@fullhuman/postcss-purgecss";
import * as htmlparser from "htmlparser2";

type PurgecssOptions = NonNullable<Parameters<typeof purgecss>[0]>;
type PurgecssExtractor = NonNullable<PurgecssOptions["extractors"]>[number];

function shouldIgnoreTag(name: string): boolean {
  if (name[0] === name[0].toUpperCase()) {
    // ignore Svelte components
    return true;
  } else if (name.startsWith("svelte:")) {
    // ignore special svelte tags
    return true;
  } else if (name == "script" || name == "style") {
    // only extract classes from html
    return true;
  }
  return false;
}

const sveltePurgecssExtractor: PurgecssExtractor["extractor"] = (
  content: string
): string[] => {
  const extracted = new Set<string>();
  const parser = new htmlparser.Parser({
    onopentag(name, attributes) {
      if (shouldIgnoreTag(name)) {
        return true;
      }

      // purgecss should keep the tag name, classes and id
      extracted.add(name);
      if (attributes.class) {
        attributes.class.split(/\s+/g)?.forEach((cls) => extracted.add(cls));
      }
      if (attributes.id) {
        extracted.add(attributes.id);
      }
    },
  });
  parser.write(content);
  parser.end();
  return Array.from(extracted);
};

export default sveltePurgecssExtractor;
