var test = require('tape');
var watchify = require('../');
var browserify = require('browserify');
var vm = require('vm');

var fs = require('fs');
var path = require('path');
var mkdirp = require('mkdirp');

var os = require('os');
var tmpbase = fs.realpathSync((os.tmpdir || os.tmpDir)());
var tmpdir = path.join(tmpbase, 'watchify-' + Math.random());

var files = {
    main: path.join(tmpdir, 'main.js'),
    beep: path.join(tmpdir, 'beep.js'),
    boop: path.join(tmpdir, 'boop.js'),
    abc: path.join(tmpdir, 'lib', 'abc.js'),
    xyz: path.join(tmpdir, 'lib', 'xyz.js')
};

mkdirp.sync(tmpdir);
mkdirp.sync(path.join(tmpdir, 'lib'));

fs.writeFileSync(files.main, [
    'var abc = require("abc");',
    'var xyz = require("xyz");',
    'var beep = require("./beep");',
    'console.log(abc + " " + xyz + " " + beep);'
].join('\n'));
fs.writeFileSync(files.beep, 'module.exports = require("./boop");');
fs.writeFileSync(files.boop, 'module.exports = require("xyz");');
fs.writeFileSync(files.abc, 'module.exports = "abc";');
fs.writeFileSync(files.xyz, 'module.exports = "xyz";');

test('properly caches exposed files', function (t) {
    t.plan(4);
    var cache = {};
    var w = watchify(browserify({
        entries: [files.main],
        basedir: tmpdir,
        cache: cache,
        packageCache: {}
    }));

    w.require('./lib/abc', {expose: 'abc'});
    w.require('./lib/xyz', {expose: 'xyz'});
    w.on('update', function () {
        w.bundle(function (err, src) {
            t.ifError(err);
            t.equal(run(src), 'ABC XYZ XYZ\n');
            w.close();
        });
    });
    w.bundle(function (err, src) {
        t.ifError(err);
        t.equal(run(src), 'abc xyz xyz\n');
        setTimeout(function () {
            // If we're incorrectly caching exposed files,
            // then "files.abc" would be re-read from disk.
            cache[files.abc].source = 'module.exports = "ABC";';
            fs.writeFileSync(files.xyz, 'module.exports = "XYZ";');
        }, 1000);
    });
});

function run (src) {
    var output = '';
    function log (msg) { output += msg + '\n' }
    vm.runInNewContext(src, { console: { log: log } });
    return output;
}
