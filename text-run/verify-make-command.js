const { cyan } = require('chalk')
const fs = require('fs')
const os = require('os')
const path = require('path')
const util = require('util')
const readFile = util.promisify(fs.readFile)

module.exports = async function ({ configuration, formatter, nodes }) {
  const expected = nodes
    .text()
    .replace(/make\s+/, '')
    .trim()
  formatter.name(`verify Make command ${cyan(expected)} exists`)
  const makefilePath = path.join(configuration.sourceDir, 'Makefile')
  const makefileContent = await readFile(makefilePath, 'utf8')
  const commands = makefileContent
    .split(os.EOL)
    .filter(lineDefinesMakeCommand)
    .map(extractMakeCommand)
  if (!commands.includes(expected)) {
    throw new Error(`Make command ${cyan(expected)} not found in ${commands}`)
  }
}

// returns whether the given line from a Makefile
// defines a Make command
function lineDefinesMakeCommand (line) {
  return makeCommandRE.test(line)
}
const makeCommandRE = /^[^ ]+:/

// returns the defined command name
// from a Makefile line that defines a Make command
function extractMakeCommand (line) {
  return line.split(':')[0]
}
