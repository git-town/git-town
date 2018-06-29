const diff = require('jsdiff-console')
const path = require('path')

module.exports = async function (activity) {
  diff(getCommand(activity), getHeadingText(activity))
}

function getHeadingText (activity) {
  return activity.nodes
    .text()
    .replace(' command', '')
    .toLowerCase()
}

function getCommand (activity) {
  return path.basename(activity.file, '.md')
}
