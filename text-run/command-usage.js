const child_process = require("child_process")
const diff = require("assert-no-diff")
const getCommand = require("./helpers/get-command.js")
const he = require("he")

module.exports = async function(activity) {
  const mdUsage = activity.nodes.text().trim()
  const cliUsage = getCliUsage(activity)
  const cliEncoded = he.encode(cliUsage, { useNamedReferences: true })
  diff.trimmedLines(mdUsage, cliEncoded)
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
