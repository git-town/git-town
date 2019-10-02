const child_process = require("child_process")
const diff = require("assert-no-diff")
const getCommand = require("./helpers/get-command.js")

module.exports = async function(activity) {
  const mdDesc = getMd(activity)
  const cliDesc = getCliDesc(activity)
  diff.trimmedLines(mdDesc, cliDesc)
}

function getMd(activity) {
  return normalize(
    activity.nodes
      .map(nodeContent)
      .map(text => text.trim())
      .join("\n")
      // Internal links should be quoted in CLI help strings,
      // while external should not
      .replace(/\n<a internal>/g, " ")
      .replace(/<\/a internal>\n/g, " ")
      .replace(/\n<\/?a external>\n/g, " ")
      .replace(/ \./g, ".")
      .replace(/ \,/g, ",")
  )
}

function nodeContent(node) {
  if (node.type === "link_open") {
    if (node.attributes.href[0] == ".") return "<a internal>"
    else return "<a external>"
  }
  if (node.type === "link_close") {
    if (node.attributes.href[0] == ".") return "</a internal>"
    else return "</a external>"
  }
  return node.content
}

function getCliDesc(activity) {
  const command = getCommand(activity.file)
  const output = child_process.execSync(`git-town help ${command}`).toString()
  const matches = output.match(/^.*\n\n([\s\S]*)\n\nUsage:\n/m)
  return normalize(matches[1].replace(/- /g, "\n").replace(/[0-9]\./g, "\n"))
}

function normalize(text) {
  return text
    .replace(/\./g, ".\n")
    .replace(/\,/g, ",\n")
    .replace(/[ ]+/g, " ")
    .replace(/\n+/g, "\n")
    .replace(/"/g, "\n")
    .replace(/:/g, "\n")
    .replace(/^\s+/gm, "")
    .replace(/\s+$/gm, "")
    .trim()
}
