
var helpers = require('./raw')
var lru     = require("lru-cache")
var cache   = lru(500)


/**
 * Meta Programming in JavaScript. Yay!
 *
 * We iterate over each file we want to memoize.
 *
 */

var fns = ['buildPriorityList', 'dataTree', 'getCurrent', 'sourceType', 'outputPath', 'outputType', 'shouldIgnore', 'findNearestLayout', 'findDefaultLayout']

fns.forEach(function(fn){
  exports[fn] = function(arg){
    var key   = fn + ':' + JSON.stringify(arguments)
    var fresh = cache.get(key)
    if(fresh) return fresh
    var hot = helpers[fn].apply(this, arguments)
    cache.set(key, hot)
    return hot
  }
})

exports.TerraformError      = helpers.TerraformError
exports.processors          = helpers.processors
exports.findFirstFile       = helpers.findFirstFile
exports.findNearestLayout   = helpers.findNearestLayout
exports.findDefaultLayout   = helpers.findDefaultLayout
exports.walkData            = helpers.walkData
exports.isTemplate          = helpers.isTemplate
exports.isStylesheet        = helpers.isStylesheet
exports.isJavaScript        = helpers.isJavaScript
exports.needsBrowserify      = helpers.needsBrowserify
