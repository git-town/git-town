const child_process = require("child_process")
const diff = require("assert-no-diff")
const getCommand = require("./helpers/get-command.js")

module.exports = async function(activity) {
  const mdDesc = getMd(activity)
  const cliDesc = getCliDesc(activity)
  diff.trimmedLines(mdDesc, cliDesc)
}

function getMd(activity) {
  let text = ""
  let isInternalLink = false
  for (node of activity.nodes) {
    switch (node.type) {
      case "text":
        text += node.content + " "
        break
      case "link_open":
        if (node.attributes.href[0] == ".") {
          isInternalLink = true
          text += '"'
        }
        break
      case "link_close":
        if (isInternalLink) {
          text += '"'
          isInternalLink = false
        }
        break
      case "code_open":
      case "code_close":
        text += '"'
        break
      case "paragraph_open":
      case "paragraph_close":
      case "list_item_open":
      case "ordered_list_open":
        text += "\n"
        break
      case "anchor_open":
      case "anchor_close":
      case "bullet_list_open":
      case "bullet_list_close":
      case "list_item_close":
        break
      default:
        throw new Error("unknown node type: " + node.type)
    }
  }
  return normalize(text.replace(/ ,/g, ",").replace(/ \./g, "."))
}

function getCliDesc(activity) {
  const command = getCommand(activity.file)
  const output = child_process.execSync(`git-town help ${command}`).toString()
  const matches = output.match(/^.*\n\n([\s\S]*)\n\nUsage:\n/m)
  return normalize(
    matches[1]
      .replace(/- /g, "\n")
      .replace(/[0-9]\./g, "\n")
      .replace(/\n\n/g, "<br>")
      .replace(/\n/g, " ")
      .replace(/<br>/g, "\n")
  )
}

function normalize(text) {
  return text
    .replace(/[ ]+/g, " ")
    .replace(/\./g, ".\n")
    .replace(/\,/g, ",\n")
    .replace(/:/g, "\n")
    .replace(/"/g, "\n")
    .replace(/\n+/g, "\n")
    .replace(/^\s+/gm, "")
    .replace(/\s+$/gm, "")
    .trim()
}
