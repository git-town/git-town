var outpipe = require('../');
var gaze = require('gaze');
var fs = require('fs');

var minimist = require('minimist');
var argv = minimist(process.argv.slice(2), {
    alias: { o: 'output' }
});

var file = argv._[0];
gaze(file, function (err, w) {
    w.on('changed', read);
});
read();

function read () {
    var r = fs.createReadStream(file);
    r.pipe(outpipe(argv.output));
}
