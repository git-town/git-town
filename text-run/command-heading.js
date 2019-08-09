const diff = require("assert-no-diff")
const getCommand = require("./helpers/get-command.js")

module.exports = async function(activity) {
  diff.trimmedLines(getCommand(activity.file), getHeadingText(activity))
}

function getHeadingText(activity) {
  return activity.nodes
    .text()
    .replace(" command", "")
    .toLowerCase()
}
