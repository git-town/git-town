const diff = require('jsdiff-console')
const getCommand = require('./helpers/get-command.js')

module.exports = async function (activity) {
  diff(getCommand(activity.file), getHeadingText(activity))
}

function getHeadingText (activity) {
  return activity.nodes
    .text()
    .replace(' command', '')
    .toLowerCase()
}
