var less           = require("less")
var TerraformError = require("../../error").TerraformError

exports.compile = function(filePath, dirs, fileContents, callback){

  var formatError = function(e){
    return new TerraformError({
      source: "Less",
      dest: "CSS",
      lineno: parseInt(e.line || -1),
      name: e.type + "Error",
      filename: filePath,
      message: e.message,
      stack: fileContents.toString()
    })
  }

  less.render(fileContents.toString(), {
      paths: dirs,            // Specify search paths for @import directives
      filename: filePath,     // Specify a filename, for better error messages
      compress: false,         // Minify CSS output
      sourceMap: true
    }, function(e, css){
      if (e) return callback(formatError(e))
      var map = css.map || ""
      return callback(null, css.css, map.toString())
  })
}

