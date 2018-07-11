const child_process = require('child_process')
const diff = require('jsdiff-console')
const getCommand = require('./helpers/get-command.js')

module.exports = async function (activity) {
  const mdSummary = activity.nodes.text().trim()
  const cliSummary = getCliDescription(activity)
  diff(mdSummary, cliSummary)
}

function getCliDescription (activity) {
  const command = getCommand(activity.file)
  const cliOutput = child_process
    .execSync(`git-town help ${command}`)
    .toString()
  const matches = cliOutput.match(/^(.*)/)
  return matches[1].trim()
}
