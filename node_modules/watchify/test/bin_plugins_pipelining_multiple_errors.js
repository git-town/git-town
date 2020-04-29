var test = require('tape');
var fs = require('fs');
var path = require('path');
var mkdirp = require('mkdirp');
var spawn = require('win-spawn');
var split = require('split');

var cmd = path.resolve(__dirname, '../bin/cmd.js');
var os = require('os');
var tmpdir = path.join((os.tmpdir || os.tmpDir)(), 'watchify-' + Math.random());

var files = {
    main: path.join(tmpdir, 'main.js'),
    plugin: path.join(tmpdir, 'plugin.js'),
    bundle: path.join(tmpdir, 'bundle.js')
};

mkdirp.sync(tmpdir);
fs.writeFileSync(files.plugin, [
    'module.exports = function(b, opts) {',
    '    b.on("file", function (file, id) {',
    '        b.pipeline.emit("error", "bad boop");',
    '        b.pipeline.emit("error", "bad boop");',
    '    });',
    '};',
].join('\n'));
fs.writeFileSync(files.main, 'boop\nbeep');

test('bin plugins pipelining multiple errors', function (t) {
    t.plan(4);
    var ps = spawn(cmd, [
        files.main,
        '-p', files.plugin, '-v',
        '-o', files.bundle
    ]);
    var lineNum = 0;
    ps.stderr.pipe(split()).on('data', function (line) {
        lineNum ++;
        if (lineNum === 1) {
            t.equal(line, 'bad boop');
        }
        if (lineNum === 2) {
            t.equal(line, 'bad boop');
            setTimeout(function() {
              fs.writeFileSync(files.main, 'beep\nboop');
            }, 1000);
        }
        if (lineNum === 3) {
            t.equal(line, 'bad boop');
        }
        if (lineNum === 4) {
            t.equal(line, 'bad boop');
            ps.kill();
        }
    });
});
