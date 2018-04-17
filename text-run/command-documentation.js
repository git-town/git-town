const child_process = require('child_process')
const diff = require('jsdiff-console')

module.exports = async function (activity) {
  const command = getCommand(activity.nodes)
  const [markdownDesc, markdownUsage, markdownFlags] = getMarkdownData(activity)
  const [gittownDesc, gittownUsage, gittownFlags] = getGittownData(command)

  diff(markdownDesc, gittownDesc)
  diff(markdownUsage, gittownUsage)
  diff(markdownFlags, gittownFlags)
}


function getCommand(nodes) {
  const commandNodes = nodes.filter(
    node => node.type === 'heading' && node.level === 1
  )
  if (commandNodes.length !== 1) throw new Error('cannot find h1')
  return commandNodes[0].content.replace(' command', '').toLowerCase()
}


function getGittownData(command) {
  const gittownOutput = child_process.execSync(`git-town help ${command}`).toString()
  const re = /(.*)\n\n([\s\S]*)\n\nUsage:\n([\s\S]*)\n\nFlags:\n([\s\S]*)/m
  const matches = gittownOutput.match(re)
  const desc = matches[1].trim()
  const usage = matches[2].trim().replace(/\n\n/g, '\n')
  const flags = matches[3].trim().replace('git-town', 'git town').replace(' [flags]', '')
  return [desc, usage, flags]
}

function getMarkdownData(activity) {
  const textBlocks = activity.nodes.filter(node => node.type === 'text')
                             .map(node => node.content)
  const desc = textBlocks[0].trim()
  const usage = textBlocks.slice(1).join('\n').trim().replace('\n\n', '\n')
  const flags = activity.searcher.tagContent('fence').trim()
  return [desc, usage, flags]
}
