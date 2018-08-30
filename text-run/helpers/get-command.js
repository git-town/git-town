const path = require('path')

module.exports = function getCommand (filename) {
  return path.basename(filename, '.md')
}
