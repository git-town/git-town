// @ts-check

const GROUP_CHARS = ["()", "<>", "[]"]

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
export function tokenize(line) {
  const tokens = []
  let token = ""
  let group = undefined
  for (const char of line) {
    if (group) {
      if (char === group[1]) {
        group = undefined
      }
      token += char
    } else {
      const nextGroup = GROUP_CHARS.find(group => group[0] === char)
      if (nextGroup) {
        group = nextGroup
        token += char
      } else if (char === " ") {
        tokens.push(token)
        token = ""
      } else {
        token += char
      }
    }
  }
  tokens.push(token)
  return tokens
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
export function extractCommand(tokens) {
  const otherTokens = [...tokens]

  const commandTokens = []
  while (otherTokens.length > 0 && otherTokens[0].match(/^[a-z]/i)) {
    commandTokens.push(otherTokens[0])
    otherTokens.shift()
  }

  const command = commandTokens.join(" ")

  return { command, otherTokens }
}
