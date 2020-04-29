var path       = require("path")
var fs         = require("fs")
var helpers    = require('../helpers')
var minify  = require('harp-minify')

/**
 * Build Processor list for javascripts.
 *
 * same as doing...
 *
 *    var processors = {
 *      "coffee" : require("./processors/coffee")
 *    }
 *
 */
var extensions = [], processors = {}
helpers.processors["js"].forEach(function(sourceType){
  extensions.push('.' + sourceType)
  processors[sourceType] = require("./processors/" + sourceType)
})

module.exports = function(root, filePath, callback){

  var srcPath = path.resolve(root, filePath)
  var ext     = path.extname(srcPath).replace(/^\./, '')
  var minifyOpts = {
    compress: false,
    mangle: false
  }

  fs.readFile(srcPath, function(err, data){

    /**
     * File not Found
     */

    if(err && err.code === 'ENOENT') return callback(null, null)

    /**
     * Read File Error
     */

    if(err) return callback(err)

    /**
     * Lookup Directories
     */

    var render = function(ext, data, cb) {
      processors[ext].compile(srcPath, data, function(err, js) {
        if (err) return cb(err)

        /**
         * Consistently minify
         */
        cb(null, minify.js(js, minifyOpts))
      })
    }

    render(ext, data, callback)

    // if(helpers.needsBrowserify(data.toString())) {
    //   var post = '', success = true

    //   var exceptionHandler = function(err) {
    //     success = false
    //     console.log(err.message)
    //     render(ext, data, callback)
    //   }

    //   process.once('uncaughtException', exceptionHandler)
    //   browserify(srcPath, {extensions: extensions}).transform(function(file) {
    //     var result = ''
    //     return through(write, end)

    //     function write(buf) {
    //       result += buf
    //     }
    //     function end() {
    //       if(success) {
    //         var that = this
    //         render(path.extname(file).replace(/^\./, '').toLowerCase(), result, function(err, data) {
    //           that.queue(data)
    //           that.queue(null)
    //         })
    //       }
    //     }
    //   }).on('error', exceptionHandler).bundle()
    //   .on('data', function(buf) {
    //     if (success) {
    //       post += buf
    //     }
    //   }).on('end', function() {
    //     if (success) {
    //       process.removeListener('uncaughtException', exceptionHandler)
    //       callback(null, minify.js(post, minifyOpts))
    //     }
    //   })
    // }
    // else {
    //   render(ext, data, callback)
    // }

  })

}
