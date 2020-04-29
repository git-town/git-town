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
    beep: path.join(tmpdir, 'beep.js'),
    boop: path.join(tmpdir, 'boop.js'),
    robot: path.join(tmpdir, 'node_modules', 'robot', 'index.js'),
    bundle: path.join(tmpdir, 'bundle.js')
};

mkdirp.sync(tmpdir);
mkdirp.sync(path.dirname(files.robot));
fs.writeFileSync(files.main, [
    'var beep = require("./beep");',
    'var boop = require("./boop");',
    'var robot = require("robot");',
    'console.log(beep + " " + boop + " " + robot);'
].join('\n'));
fs.writeFileSync(files.beep, 'module.exports = "beep";');
fs.writeFileSync(files.boop, 'module.exports = "boop";');
fs.writeFileSync(files.robot, 'module.exports = "robot";');

test('api ignore watch', function (t) {
    t.plan(4);
    var ps = spawn(cmd, [
        files.main,
        '--ignore-watch', '**/be*.js',
        '-o', files.bundle,
        '-v'
    ]);
    var lineNum = 0;
    ps.stderr.pipe(split()).on('data', function (line) {
        lineNum ++;
        if (lineNum === 1) {
            run(files.bundle, function (err, output) {
                t.ifError(err);
                t.equal(output, 'beep boop robot\n');
                fs.writeFileSync(files.beep, 'module.exports = "BEEP";');
                fs.writeFileSync(files.boop, 'module.exports = "BOOP";');
                fs.writeFileSync(files.robot, 'module.exports = "ROBOT";');
            });
        }
        else if (lineNum === 2) {
            run(files.bundle, function (err, output) {
                t.ifError(err);
                t.equal(output, 'beep BOOP ROBOT\n');
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
