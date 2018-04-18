const child_process = require('child_process')
const diff = require('jsdiff-console')
const path = require('path')

module.exports = async function (activity) {
  const markdownDesc = activity.searcher.tagContent('fence')
  const gittownDesc = getGittownDescription(activity)
  diff(markdownDesc, gittownDesc)
}


function getCommand(activity) {
  return path.basename(activity.filename, '.md')
}


function getGittownDescription(activity) {
  const command = getCommand(activity)
  const gittownOutput = child_process.execSync(`git-town help ${command}`).toString()
  const matches = gittownOutput.match(/\nUsage:\n(.*)/)
  return matches[1].trim().replace(' [flags]', '').replace('git-town', 'git town')
}
