const diff = require('jsdiff-console')
const path = require('path')

module.exports = async function (activity) {
  diff(getCommand(activity), getHeadingText(activity))
}


function getHeadingText(activity) {
  return activity.searcher.tagContent('heading')
                          .replace(' command', '')
                          .toLowerCase()
}


function getCommand(activity) {
  return path.basename(activity.filename, '.md')
}
