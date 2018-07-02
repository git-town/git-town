const child_process = require('child_process')
const diff = require('jsdiff-console')
const getCommand = require('./helpers/get-command.js')

module.exports = async function (activity) {
  const markdownDesc = activity.nodes.text().trim()
  const gittownDesc = getGittownDescription(activity)
  diff(markdownDesc, gittownDesc)
}

function getGittownDescription (activity) {
  const command = getCommand(activity.file)
  const gittownOutput = child_process
    .execSync(`git-town help ${command}`)
    .toString()
  const matches = gittownOutput.match(/\nUsage:\n(.*)/)
  return matches[1]
    .trim()
    .replace(' [flags]', '')
    .replace('git-town', 'git town')
}
