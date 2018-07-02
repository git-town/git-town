const child_process = require('child_process')
const diff = require('jsdiff-console')
const getCommand = require('./helpers/get-command.js')

module.exports = async function (activity) {
  const markdownDesc = activity.nodes.text().trim()
  const gittownDesc = getGittownFlags(activity)
  diff([markdownDesc], gittownDesc)
}

function getGittownFlags (activity) {
  return child_process
    .execSync(`git-town help ${getCommand(activity.file)}`)
    .toString()
    .match(/\nFlags:\n([\s\S]*)\nGlobal Flags:\n/)[1]
    .split('\n')
    .filter(line => line)
    .filter(line => !line.includes('help'))
    .map(line => line.trim())
}
