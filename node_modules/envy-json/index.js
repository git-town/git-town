
function injectEnv(key, value) {
  var dollar = new RegExp(/^\$/)
  if (typeof(value) == 'string' && value.match(dollar)) {
    return process.env[value.replace(dollar, '')]
  }
  return value
}

module.exports = function(json) {
  if (typeof(json) === "object") {
    return JSON.parse(JSON.stringify(json, injectEnv))
  } else if ((typeof(json) === "string") && json.match(/\.json$/i)) {
    return JSON.parse(JSON.stringify(JSON.parse(require('fs').readFileSync(json)), injectEnv))
  } else {
    return JSON.parse(JSON.stringify(JSON.parse(json), injectEnv))
  }
}