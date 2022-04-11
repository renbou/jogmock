/**
 * RegExp-based extractor patterns taken from TailwindCss.
 * This is contained in tailwindcss/src/lib/defaultExtractor.js
 */
const PATTERNS = [
  /(?:\['([^'\s]+[^<>"'`\s:\\])')/.source, // ['text-lg' -> text-lg
  /(?:\["([^"\s]+[^<>"'`\s:\\])")/.source, // ["text-lg" -> text-lg
  /(?:\[`([^`\s]+[^<>"'`\s:\\])`)/.source, // [`text-lg` -> text-lg
  /([^${(<>"'`\s]*\[\w*'[^"`\s]*'?\])/.source, // font-['some_font',sans-serif]
  /([^${(<>"'`\s]*\[\w*"[^'`\s]*"?\])/.source, // font-["some_font",sans-serif]
  /([^<>"'`\s]*\[\w*\('[^"'`\s]*'\)\])/.source, // bg-[url('...')]
  /([^<>"'`\s]*\[\w*\("[^"'`\s]*"\)\])/.source, // bg-[url("...")]
  /([^<>"'`\s]*\[\w*\('[^"`\s]*'\)\])/.source, // bg-[url('...'),url('...')]
  /([^<>"'`\s]*\[\w*\("[^'`\s]*"\)\])/.source, // bg-[url("..."),url("...")]
  /([^<>"'`\s]*\[[^<>"'`\s]*\('[^"`\s]*'\)+\])/.source, // h-[calc(100%-theme('spacing.1'))]
  /([^<>"'`\s]*\[[^<>"'`\s]*\("[^'`\s]*"\)+\])/.source, // h-[calc(100%-theme("spacing.1"))]
  /([^${(<>"'`\s]*\['[^"'`\s]*'\])/.source, // `content-['hello']` but not `content-['hello']']`
  /([^${(<>"'`\s]*\["[^"'`\s]*"\])/.source, // `content-["hello"]` but not `content-["hello"]"]`
  /([^<>"'`\s]*\[[^<>"'`\s]*:[^\]\s]*\])/.source, // `[attr:value]`
  /([^<>"'`\s]*\[[^<>"'`\s]*:'[^"'`\s]*'\])/.source, // `[content:'hello']` but not `[content:"hello"]`
  /([^<>"'`\s]*\[[^<>"'`\s]*:"[^"'`\s]*"\])/.source, // `[content:"hello"]` but not `[content:'hello']`
  /([^<>"'`\s]*\[[^"'`\s]+\][^<>"'`\s]*)/.source, // `fill-[#bada55]`, `fill-[#bada55]/50`
  /([^"'`\s]*[^<>"'`\s:\\])/.source, //  `<sm:underline`, `md>:font-bold`
  /([^<>"'`\s]*[^"'`\s:\\])/.source, //  `px-1.5`, `uppercase` but not `uppercase:`
].join("|");

const BROAD_MATCH_GLOBAL_REGEXP = new RegExp(PATTERNS, "g");
const INNER_MATCH_GLOBAL_REGEXP =
  /[^<>"'`\s.(){}[\]#=%$]*[^<>"'`\s.(){}[\]#=%:$]/g;

export function defaultExtractor(content: string) {
  let broadMatches = content.matchAll(BROAD_MATCH_GLOBAL_REGEXP);
  let innerMatches = content.match(INNER_MATCH_GLOBAL_REGEXP) || [];
  let results = [...broadMatches, ...innerMatches]
    .flat()
    .filter((v) => v !== undefined);

  return results;
}

/**
 * Svelte-specific extractor which parses .svelte and takes class names, ids, etc.
 * Should not be used as a simpler RegExp solution works just fine.
 */
// function shouldIgnoreTag(name: string): boolean {
//   if (name[0] === name[0].toUpperCase()) {
//     // ignore Svelte components
//     return true;
//   } else if (name.startsWith("svelte:")) {
//     // ignore special svelte tags
//     return true;
//   } else if (name == "script" || name == "style") {
//     // only extract classes from html
//     return true;
//   }
//   return false;
// }

// const sveltePurgecssExtractor: PurgecssExtractor["extractor"] = (
//   content: string
// ): string[] => {
//   const extracted = new Set<string>();
//   const parser = new htmlparser.Parser({
//     onopentag(name, attributes) {
//       if (shouldIgnoreTag(name)) {
//         return true;
//       }

//       // purgecss should keep the tag name, classes and id
//       extracted.add(name);
//       if (attributes.class) {
//         attributes.class.split(/\s+/g)?.forEach((cls) => extracted.add(cls));
//       }
//       if (attributes.id) {
//         extracted.add(attributes.id);
//       }
//     },
//   });
//   parser.write(content);
//   parser.end();
//   return Array.from(extracted);
// };

// export default sveltePurgecssExtractor;
