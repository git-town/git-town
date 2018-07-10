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
    .map(nodeContent)
    .map(text => text.trim())
    .join('\n')
    .replace(/\./g, '.\n')
    .replace(/\,/g, ',\n')
    .replace(/:/g, '\n')
    .replace(/"/g, '\n')
    .replace(/\n<\/?a>\n/g, ' ')
    .replace(/[ ]+/g, ' ')
    .replace(/\n+/g, '\n')
    .trim()
}

function nodeContent (node) {
  if (node.type === 'link_open') return '<a>'
  if (node.type === 'link_close') return '</a>'
  return node.content
}

function getGittownUsage (activity) {
  const command = getCommand(activity.file)
  const output = child_process.execSync(`git-town help ${command}`).toString()
  const matches = output.match(/^.*\n\n([\s\S]*)\n\nUsage:\n/m)
  return matches[1]
    .replace(/[ ]+/g, ' ')
    .replace(/\./g, '.\n')
    .replace(/\,/g, ',\n')
    .replace(/:/g, '\n')
    .replace(/- /g, '\n')
    .replace(/[0-9]\./g, '\n')
    .replace(/"/g, '\n')
    .replace(/^\s+/gm, '\n')
    .replace(/\s+$/gm, '\n')
    .replace(/\n+/g, '\n')
    .trim()
}
