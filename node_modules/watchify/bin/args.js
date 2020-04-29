var fromArgs = require('browserify/bin/args');
var watchify = require('../');
var defined = require('defined');
var xtend = require('xtend');

module.exports = function (args) {
    var b = fromArgs(args, watchify.args);

    var opts = {};
    var ignoreWatch = defined(b.argv['ignore-watch'], b.argv.iw);
    if (ignoreWatch) {
        opts.ignoreWatch = ignoreWatch;
    }

    return watchify(b, xtend(opts, b.argv));
};
