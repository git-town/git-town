const child_process = require('child_process')
const diff = require('jsdiff-console')
const getCommand = require('./helpers/get-command.js')

module.exports = async function (activity) {
  const mdUsage = getMd(activity)
  const gittownUsage = getGittownUsage(activity)
  diff(mdUsage, gittownUsage)
}

function getMd (activity) {
  return activity.nodes
    .map(node => node.content)
    .filter(node => node)
    .join('\n')
    .replace(/[ ]+/g, ' ')
    .replace(/\./g, '.\n')
    .replace(/\,/g, ',\n')
    .replace(/:/g, ':\n')
    .replace(/"/g, '\n')
    .replace(/^\s*/gm, '')
    .replace(/\s*$/gm, '')
    .trim()
}

function getGittownUsage (activity) {
  const command = getCommand(activity.file)
  const output = child_process.execSync(`git-town help ${command}`).toString()
  const matches = output.match(/^.*\n\n([\s\S]*)\n\nUsage:\n/m)
  return matches[1]
    .replace(/[ ]+/g, ' ')
    .replace(/\./g, '.\n')
    .replace(/\,/g, ',\n')
    .replace(/:/g, ':\n')
    .replace(/- /g, '\n')
    .replace(/[0-9]\./g, '\n')
    .replace(/"/g, '\n')
    .replace(/^\s+/gm, '\n')
    .replace(/\s+$/gm, '\n')
    .replace(/\n+/g, '\n')
    .trim()
}
