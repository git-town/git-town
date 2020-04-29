// @see http://nodejs.org/api/all.html#all_require_extensions
var fs = require('fs')
  , yaml = require('js-yaml');

require.extensions['.yaml'] = 
require.extensions['.yml'] = function(module, filename) {
  var content = fs.readFileSync(filename, 'utf8');
  // Parse the file content and give to module.exports
  content = yaml.load(content);
  module.exports = content;
};
