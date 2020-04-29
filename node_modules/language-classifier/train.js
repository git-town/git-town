
/**
 * Module dependencies.
 */

var fs = require('fs')
  , Bayesian = require('classifier').Bayesian
  , path = require('path')
  , read = fs.readFileSync
  , write = fs.writeFileSync
  , readdir = fs.readdirSync;

var classifier = new Bayesian;

console.log();
readdir('training-set').forEach(function(file){
  file = path.join('training-set', file);
  var ext = path.extname(file);
  var lang = path.basename(file, ext);
  console.log('  %s : %s', 'train', lang, ext);
  var str = read(file, 'utf8');
  classifier.train(str, lang);
});

var json = JSON.stringify(classifier);
console.log('  write : memory.json');
write('memory.json', json);
console.log();
