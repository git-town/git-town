import * as textRunner from "text-runner"
import { SummarySection } from "./summary.ts"

/** Document contains the AST for an entire webpage describing a Git Town command */
export class Document {
  nodes: textRunner.ast.NodeList

  constructor(nodes: textRunner.ast.NodeList) {
    this.nodes = nodes
  }

  /** provides the text of the ```command-summary block at the beginning of this page */
  summarySection(): SummarySection {
    const fences = this.nodes.nodesOfTypes("fence")
    if (fences.length === 0) {
      throw new Error("no fenced blocks found")
    }
    // the first fenced block contains the summary
    const summaryBlock = fences[0]
    const summaryNodes = this.nodes.nodesFor(summaryBlock)
    return new SummarySection(summaryNodes.text())
  }

  /** provides the arguments documented in the "## Options" section of this page */
  argsInOptions(): string[][] {
    let result: string[][] = []
    let insideOptions = false
    for (const node of this.nodes) {
      if (isH2(node)) {
        if (insideOptions) {
          // here we run into the next h2 heading after options --> done parsing options
          return result
        }
        if (this.isOptionsHeading(node)) {
          insideOptions = true
        }
        continue
      }
      if (insideOptions) {
        if (isFlagHeading(node)) {
          const flagNodes = this.nodes.nodesFor(node)
          result.push(texts(flagNodes))
        }
      }
    }
    return result
  }

  /** indicates whether the given node is the "## Options" heading */
  isOptionsHeading(node: textRunner.ast.Node): boolean {
    const nodes = this.nodes.nodesFor(node)
    const text = nodes.text()
    return text === "Options"
  }
}

function isFlagHeading(node: textRunner.ast.Node): boolean {
  return node.type === "h4_open"
}

function isH2(node: textRunner.ast.Node): boolean {
  return node.type === "h2_open"
}

function texts(nodes: textRunner.ast.NodeList): string[] {
  let result: string[] = []
  for (const node of nodes) {
    if (node.type === "text") {
      result.push(node.content)
    }
  }
  return standardizeArgument(removeNegatedFlag(result))
}

export function removeNegatedFlag(flags: string[]): string[] {
  if (flags.length < 2) {
    return flags
  }
  return flags.filter(flag => !flag.startsWith("--no-"))
}

export function standardizeArgument(texts: string[]): string[] {
  const result: string[] = []
  for (const text of texts) {
    if (text.startsWith("--")) {
      const parts = text.split(" ")
      if (parts.length > 1) {
        parts[1] = "string"
      }
      result.push(parts.join(" "))
    } else if (text.startsWith("-")) {
      const parts = text.split(" ")
      result.push(parts[0])
    }
  }
  return result
}
