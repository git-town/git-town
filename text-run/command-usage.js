const child_process = require("child_process")
const diff = require("assert-no-diff")
const getCommand = require("./helpers/get-command.js")

module.exports = async function(activity) {
  const mdUsage = activity.nodes
    .text()
    .trim()
    .replace(/&lt;/g, "<")
    .replace(/&gt;/g, ">")
  const cliUsage = getCliUsage(activity)
  diff.trimmedLines(mdUsage, cliUsage)
}

function getCliUsage(activity) {
  const command = getCommand(activity.file)
  const cliOutput = child_process
    .execSync(`git-town help ${command}`)
    .toString()
  const matches = cliOutput.match(/\nUsage:\n([\s\S]*?)\n\n/)
  return matches[1]
    .trim()
    .replace(" [flags]", "")
    .replace(/git-town/g, "git town")
    .replace(/^\s+/gm, "")
}
