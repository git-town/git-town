const child_process = require('child_process')
const diff = require('jsdiff-console')
const getCommand = require('./helpers/get-command.js')

module.exports = async function (activity) {
  const markdownUsage = activity.nodes
    .text()
    .replace(/\./g, '.\n')
    .replace(/\s+/, ' ')
  const gittownUsage = getGittownUsage(activity)
  diff(markdownUsage, gittownUsage)
}

function getGittownUsage (activity) {
  const command = getCommand(activity.file)
  const output = child_process.execSync(`git-town help ${command}`).toString()
  const matches = output.match(/^.*\n\n([\s\S]*)\n\nUsage:\n/m)
  return matches[1]
    .trim()
    .replace(/\n/g, '')
    .replace(/\s+/g, ' ')
    .replace(/\./g, '.\n')
}
