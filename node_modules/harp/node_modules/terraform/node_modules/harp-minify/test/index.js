
var assert = require('assert');
var exec = require('child_process').exec;
var fs = require('fs');
var minify = require('..');
var path = require('path');

/**
 * Languages.
 */

var langs = [
  'css',
  'js',
  'html'
];

/**
 * Node API.
 */

describe('minify', function () {
  langs.forEach(function (lang) {
    var min = fixture('min.' + lang);
    var max = fixture('max.' + lang);

    it(lang, function () {
      assert.equal(min, minify[lang](max));
    });

    it(lang + ' (unknown)', function() {
      assert.equal(min, minify(max));
    });
  });

  it('should pass options to uglify js', function(){
    var js = '(function(){ var analytics = window.analytics = []; analytics.push(1) })();';
    var opts = { mangle: { except: ['analytics'] }};
    assert(!~minify(js).indexOf('var analytics'), 'expected "analytics" to be mangled');
    assert(~minify(js, opts).indexOf('var analytics'), 'expected "analytics" to not be mangled');
  });
});

/**
 * CLI.
 */

describe('cli', function () {
  afterEach(function (done) {
    exec('rm -rf output', done);
  });

  langs.forEach(function (lang) {
    it(lang, function (done) {
      var min = fixture('min.' + lang);
      exec('bin/minify test/fixtures/max.' + lang, function (err, stdout) {
        if (err) return done(err);
        assert(stdout == min);
        done();
      });
    });
  });
});

/**
 * Read a fixture by name.
 *
 * @param {String} name
 */

function fixture (name) {
  var filename = path.resolve(__dirname, 'fixtures', name);
  return fs.readFileSync(filename, 'utf8');
}
