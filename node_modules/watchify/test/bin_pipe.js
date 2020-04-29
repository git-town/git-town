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
    bundle: path.join(tmpdir, 'bundle.js')
};

mkdirp.sync(tmpdir);
fs.writeFileSync(files.main, 'console.log(num * 2)');

test('bin with pipe', function (t) {
    t.plan(5);
    var ps = spawn(cmd, [
        files.main,
        '-o', 'uglifyjs - --enclose 11:num > ' + files.bundle,
        '-v'
    ]);
    var lineNum = 0;
    ps.stderr.pipe(split()).on('data', function (line) {
        lineNum ++;
        if (lineNum === 1) {
            run(files.bundle, function (err, output) {
                t.ifError(err);
                t.equal(output, '22\n');
                fs.writeFile(files.main, 'console.log(num * 3)', t.ifError);
            });
        }
        else if (lineNum === 2) {
            run(files.bundle, function (err, output) {
                t.ifError(err);
                t.equal(output, '33\n');
                ps.kill();
            });
        }
    });
});

function run (file, cb) {
    var ps = spawn(process.execPath, [ file ]);
    var data = [];
    ps.stdout.on('data', function (buf) { data.push(buf) });
    ps.stdout.on('end', function () {
        cb(null, Buffer.concat(data).toString('utf8'));
    });
    ps.on('error', cb);
    return ps;
}
