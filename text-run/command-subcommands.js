const child_process = require("child_process")
const diff = require("assert-no-diff")
const getCommand = require("./helpers/get-command.js")

module.exports = async function(activity) {
  const mdCommands = await getMdCommands(activity.nodes)
  const cliCommands = getCliCommands(activity)
  diff.trimmedLines(mdCommands.join("\n"), cliCommands.join("\n"))
}

function getCliCommands(activity) {
  const result = []
  const command = getCommand(activity.file)
  const cliOutput = child_process
    .execSync(`git-town help ${command}`)
    .toString()
  const matches = cliOutput.match(/\nAvailable Commands:\n([\s\S]*?)\n\n/)
  const text = matches[1]
  for (const line of text.split("\n")) {
    const words = line.trim().split(/\s+/)
    const command = words[0]
    const desc = words.slice(1).join(" ")
    result.push([command, desc])
  }
  return result
}

async function getMdCommands(nodes) {
  const result = []
  const table = nodes.getNodeOfTypes("table_open")
  const tableNodes = nodes.getNodesFor(table)
  for (const tableRow of tableNodes.getNodesOfTypes("table_row_open")) {
    const rowNodes = tableNodes.getNodesFor(tableRow)
    const commandName = rowNodes.getNodeOfTypes("table_heading")
    const commandDesc = rowNodes.getNodeOfTypes("table_cell")
    result.push([commandName.content, commandDesc.content])
  }
  return result
}
