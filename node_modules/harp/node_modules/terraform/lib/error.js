
function TerraformError(error) {
  this.source   = error.source
  this.dest     = error.dest
  this.filename = error.filename
  this.lineno   = error.lineno
  this.name     = error.name
  this.message  = error.message
  this.stack    = error.stack
}

TerraformError.prototype = Error.prototype;

exports.TerraformError = TerraformError