var ejs = require('ejs')
var TerraformError = require("../../error").TerraformError

module.exports = function(fileContents, options){

  return {
    compile: function(){
      return ejs.compile(fileContents.toString(), options)
    },

    parseError: function(error){
      var arr = error.message.split("\n")
      var path_arr = arr[0].split(":")

      error.lineno   = parseInt(error.lineno || path_arr[path_arr.length -1] || -1)
      error.message  = arr[arr.length - 1]
      error.name     = error.name
      error.source   = "EJS"
      error.dest     = "HTML"
      error.filename = error.path || options.filename
      error.stack    = fileContents.toString()

      return new TerraformError(error)
    }
  }

}