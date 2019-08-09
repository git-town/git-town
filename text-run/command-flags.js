const child_process = require("child_process")
const diff = require("assert-no-diff")
const getCommand = require("./helpers/get-command.js")

module.exports = async function(activity) {
  const mdFlags = getMdFlags(activity)
  const cliFlags = getCliFlags(activity)
  diff.trimmedLines(mdFlags, cliFlags)
}

function getMdFlags(activity) {
  return activity.nodes.text().trim()
}

function getCliFlags(activity) {
  return child_process
    .execSync(`git-town help ${getCommand(activity.file)}`)
    .toString()
    .match(/\nFlags:\n([\s\S]*)\nGlobal Flags:\n/)[1]
    .split("\n")
    .filter(line => line)
    .filter(line => !line.includes("help"))
    .map(line => line.trim())
    .join("\n")
}
