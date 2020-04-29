# outpipe

write output to a file through shell commands

# purpose

Suppose you have a tool like
[watchify](https://npmjs.com/package/watchify)
or [factor-bundle](https://npmjs.com/package/factor-bundle)
that write to multiple files or write to the same file more than once.

If you want to pipe the output of these tools to other programs, such as
minification with the `uglify` command, it's very difficult! You might need to
use the tool's API or use a separate command to watch for changes to the output
files. Ick.

You don't get the elegance of something like:

``` sh
$ browserify main.js | uglifyjs -cm | gzip > bundle.js.gz
```

Until now! With this library and a hypothetical version of watchify, you could
do:

``` sh
$ watchify main.js -dv -o 'uglifyjs -cm | gzip > bundle.js.gz'
```

# example

Here's a small watcher program that will just copy input files to a destination,
but transforms can be applied along the way with shell pipes and redirects.

``` js
var outpipe = require('outpipe');
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
```

We can run the program with a single output file:

``` sh
$ node watch.js input/x.js -o output/hmm.js
```

which just copies `x.js` to `output/hmm.js` whenever `x.js` changes.

We could also run a minification step using the `uglify` command:

``` sh
$ node watch.js input/x.js -o 'uglifyjs -cm > output/wow.js'
```

or we can just print the size of the minified and gzipped output to stdout:

``` sh
$ node watch.js input/x.js -o 'uglifyjs -cm | gzip | wc -c'
123
```

or we could write that size to a file:

``` sh
$ node watch.js input/x.js -o 'uglifyjs -cm | gzip | wc -c > size.txt'
```

# methods

``` js
var outpipe = require('outpipe')
```

## var w = outpipe(cmd, opts={})

Return a writable stream `w` that will pipe output to the command string `cmd`.

If `cmd` has no operators (`|` or `>`), it will write to a file.

Otherwise, each command between pipes will be executed and output is written to
a file if `>` is given.

`opts` can be:

* `opts.env` - an object mapping environment variables to their values or a
`function (key) {}` that returns the values.

stdout and stderr are forwarded to process.stdout and process.stderr if
unhandled in the command.

# install

With [npm](https://npmjs.org) do:

```
npm install outpipe
```

# license

MIT
