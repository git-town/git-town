// @ts-check
/// <reference types="node" />

/**
 * @typedef {import("./types").Book} Book
 * @typedef {import("./types").BookItem} BookItem
 * @typedef {import("./types").Chapter} Chapter
 *
 * This file is an mdBook preprocessor that modifies the book JSON data.
 * We use it to process code blocks with language "command-summary".
 *
 * @see processCommandSummary
 *
 * mdBook preprocessor documentation:
 * https://rust-lang.github.io/mdBook/for_developers/preprocessors.html
 */

if (process.argv.length > 2) {
  handleSupports();
} else {
  await handlePreprocess();
}

function handleSupports() {
  if (process.argv[2] === "supports" && process.argv[3] === "html") {
    process.exit(0);
  }
  process.exit(1);
}

/**
 * An mdBook preprocessor receives JSON data `[context, book]` from stdin,
 * modifies `book`, and writes `book` back to stdout.
 *
 * Example input:
 * https://rust-lang.github.io/mdBook/for_developers/preprocessors.html#:~:text=mod%20test
 */
async function handlePreprocess() {
  // Read from stdin
  let stdin = "";
  for await (const chunk of process.stdin) {
    stdin += chunk;
  }

  // We don't care about the context. Only process the book.
  const [, book] = JSON.parse(stdin);

  processBook(book);

  // Write to stdout
  const output = JSON.stringify(book);
  process.stdout.write(output + "\n");
}

/**
 * @param {Book} book
 */
function processBook(book) {
  for (const bookItem of book.sections) {
    processBookItem(bookItem);
  }
}

/**
 * @param {BookItem} bookItem
 */
function processBookItem(bookItem) {
  // bookItem is { "Chapter": Chapter }, "Separator", or { "PartTitle": string }
  if (bookItem === "Separator") {
    return;
  }
  // bookItem is { "Chapter": Chapter } or { "PartTitle": string }
  if ("PartTitle" in bookItem) {
    return;
  }
  // bookItem is { "Chapter": Chapter }

  processChapter(bookItem.Chapter);
}

/**
 * @param {Chapter} chapter
 */
function processChapter(chapter) {
  for (const subItem of chapter.sub_items) {
    processBookItem(subItem);
  }

  chapter.content = processContent(chapter.content);
}

/**
 * @param {string} content
 * @returns {string}
 */
function processContent(content) {
  return content.replaceAll(/```command-summary\n([\s\S]*?)\n```/g, (_, code) => {
    return processCommandSummary(code);
  });
}

/**
 * This function processes code blocks with language "command-summary".
 *
 * For example:
 *
 * ```command-summary
 * git town append [-p | --prototype] [-v | --verbose] <branch-name>
 * ```
 *
 * will become:
 *
 * <pre><code><div class="gt-command-summary" style="padding-left: 16ch; text-indent: -16ch"
 *   ><span class="gt-command">git town append</span
 *   > <span>[-p | --prototype]</span
 *   > <span>[-v | --verbose]</span
 *   > <span>&lt;branch-name&gt;</span
 *   ></div></code></pre>
 *
 * `padding-left` and `text-indent` are set based on the length of the command.
 * They align the arguments with the command. Other styles are applied in
 * head.hbs. The above example should render as:
 *
 * ┌───────────────────────────────────────────────────────────────────┐
 * │ git town append [-p | --prototype] [-v | --verbose] <branch-name> │
 * └───────────────────────────────────────────────────────────────────┘
 * or
 * ┌─────────────────────────────────────────────────────┐
 * │ git town append [-p | --prototype] [-v | --verbose] │
 * │                 <branch-name>                       │
 * └─────────────────────────────────────────────────────┘
 * or
 * ┌───────────────────────────────────────┐
 * │ git town append [-p | --prototype]    │
 * │                 [-v | --verbose]      │
 * │                 <branch-name>         │
 * └───────────────────────────────────────┘
 *
 * @param {string} code
 * @returns {string}
 */
function processCommandSummary(code) {
  return `<pre><code>${
    code
      .split("\n")
      .map(line => {
        const tokens = tokenize(line);
        const { command, otherTokens } = extractCommand(tokens);

        const indent = command.length + 1;
        return `<div class="gt-command-summary" style="padding-left: ${indent}ch; text-indent: -${indent}ch"><span class="gt-command">${command}</span> ${
          otherTokens
            .map(token => `<span>${token.replaceAll("<", "&lt;").replaceAll(">", "&gt;")}</span>`)
            .join(" ")
        }</div>`;
      })
      .join("")
  }</pre></code>`;
}

/**
 * This function tokenizes a line of text into strings that should be kept
 * together when wrapping text. For example, the text "[-p | --prototype]"
 * should be a single token.
 *
 * @example
 * tokenize("git town append [-p | --prototype] <branch-name>")
 * // => ["git", "town", "append", "[-p | --prototype]", "<branch-name>"]
 *
 * @param {string} line
 * @returns {string[]}
 */
function tokenize(line) {
  const GROUP_CHARS = ["()", "<>", "[]"];

  const tokens = [];
  let token = "";
  let group = undefined;
  for (const char of line) {
    if (group) {
      if (char === group[1]) {
        group = undefined;
      }
      token += char;
    } else {
      const nextGroup = GROUP_CHARS.find(group => group[0] === char);
      if (nextGroup) {
        group = nextGroup;
        token += char;
      } else if (char === " ") {
        tokens.push(token);
        token = "";
      } else {
        token += char;
      }
    }
  }
  tokens.push(token);
  return tokens;
}

/**
 * This function extracts the command and other tokens from a line of text.
 *
 * @example
 * extractCommand(["git", "town", "append", "[-p | --prototype]", "<branch-name>"])
 * // => { command: "git town append", otherTokens: ["[-p | --prototype]", "<branch-name>"] }
 *
 * @param {string[]} tokens
 * @returns {{ command: string, otherTokens: string[] }}
 */
function extractCommand(tokens) {
  const otherTokens = [...tokens];

  const commandTokens = [];
  while (otherTokens.length > 0 && otherTokens[0].match(/^[a-z]/i)) {
    commandTokens.push(otherTokens[0]);
    otherTokens.shift();
  }

  const command = commandTokens.join(" ");

  return { command, otherTokens };
}

// Make TypeScript think this file is a module
export {};
