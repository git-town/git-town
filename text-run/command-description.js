const child_process = require('child_process')
const diff = require('jsdiff-console')
const getCommand = require('./helpers/get-command.js')

module.exports = async function (activity) {
  const mdDesc = getMd(activity)
  const cliDesc = getCliDesc(activity)
  diff(mdDesc, cliDesc)
}

function getMd (activity) {
  return normalize(
    activity.nodes
      .map(nodeContent)
      .map(text => text.trim())
      .join('\n')
      .replace(/\n<\/?a>\n/g, ' ')
      .replace(/ \./g, '.')
      .replace(/ \,/g, ',')
  )
}

function nodeContent (node) {
  if (node.type === 'link_open') return '<a>'
  if (node.type === 'link_close') return '</a>'
  return node.content
}

function getCliDesc (activity) {
  const command = getCommand(activity.file)
  const output = child_process.execSync(`git-town help ${command}`).toString()
  const matches = output.match(/^.*\n\n([\s\S]*)\n\nUsage:\n/m)
  return normalize(matches[1].replace(/- /g, '\n').replace(/[0-9]\./g, '\n'))
}

function normalize (text) {
  return text
    .replace(/\./g, '.\n')
    .replace(/\,/g, ',\n')
    .replace(/[ ]+/g, ' ')
    .replace(/\n+/g, '\n')
    .replace(/"/g, '\n')
    .replace(/:/g, '\n')
    .replace(/^\s+/gm, '')
    .replace(/\s+$/gm, '')
    .trim()
}
