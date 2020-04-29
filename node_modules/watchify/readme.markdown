# watchify

watch mode for [browserify](https://github.com/substack/node-browserify) builds

[![build status](https://secure.travis-ci.org/substack/watchify.svg?branch=master)](http://travis-ci.org/substack/watchify)

Update any source file and your browserify bundle will be recompiled on the
spot.

# example

```
$ watchify main.js -o static/bundle.js
```

Now as you update files, `static/bundle.js` will be automatically
incrementally rebuilt on the fly.

The `-o` option can be a file or a shell command (not available on Windows)
that receives piped input:

``` sh
watchify main.js -o 'exorcist static/bundle.js.map > static/bundle.js' -d
```

``` sh
watchify main.js -o 'uglifyjs -cm > static/bundle.min.js'
```

You can use `-v` to get more verbose output to show when a file was written and how long the bundling took (in seconds):

```
$ watchify browser.js -d -o static/bundle.js -v
610598 bytes written to static/bundle.js (0.23 seconds) at 8:31:25 PM
610606 bytes written to static/bundle.js (0.10 seconds) at 8:45:59 PM
610597 bytes written to static/bundle.js (0.14 seconds) at 8:46:02 PM
610606 bytes written to static/bundle.js (0.08 seconds) at 8:50:13 PM
610597 bytes written to static/bundle.js (0.08 seconds) at 8:58:16 PM
610597 bytes written to static/bundle.js (0.19 seconds) at 9:10:45 PM
```

# usage

Use `watchify` with all the same options as `browserify` except that `-o` (or
`--outfile`) is mandatory. Additionally, there are also:

```
Standard Options:

  --outfile=FILE, -o FILE

    This option is required. Write the browserify bundle to this file. If
    the file contains the operators `|` or `>`, it will be treated as a
    shell command, and the output will be piped to it.

  --verbose, -v                     [default: false]

    Show when a file was written and how long the bundling took (in
    seconds).

  --version

    Show the watchify and browserify versions with their module paths.
```

```
Advanced Options:

  --delay                           [default: 100]

    Amount of time in milliseconds to wait before emitting an "update"
    event after a change.

  --ignore-watch=GLOB, --iw GLOB    [default: false]

    Ignore monitoring files for changes that match the pattern. Omitting
    the pattern will default to "**/node_modules/**".

  --poll=INTERVAL                   [default: false]

    Use polling to monitor for changes. Omitting the interval will default
    to 100ms. This option is useful if you're watching an NFS volume.
```

# methods

``` js
var watchify = require('watchify');
```

## watchify(b, opts)

watchify is a browserify [plugin](https://github.com/substack/node-browserify#bpluginplugin-opts), so it can be applied like any other plugin.
However, when creating the browserify instance `b`, **you MUST set the `cache`
and `packageCache` properties**:

``` js
var b = browserify({ cache: {}, packageCache: {} });
b.plugin(watchify);
```

```js
var b = browserify({
  cache: {},
  packageCache: {},
  plugin: [watchify]
});
```

**By default, watchify doesn't display any output, see [events](https://github.com/substack/watchify#events) for more info.**

`b` continues to behave like a browserify instance except that it caches file
contents and emits an `'update'` event when a file changes. You should call
`b.bundle()` after the `'update'` event fires to generate a new bundle.
Calling `b.bundle()` extra times past the first time will be much faster due
to caching.

**Important:** Watchify will not emit `'update'` events until you've called
`b.bundle()` once and completely drained the stream it returns.

```js
var fs = require('fs');
var browserify = require('browserify');
var watchify = require('watchify');

var b = browserify({
  entries: ['path/to/entry.js'],
  cache: {},
  packageCache: {},
  plugin: [watchify]
});

b.on('update', bundle);
bundle();

function bundle() {
  b.bundle()
    .on('error', console.error)
    .pipe(fs.createWriteStream('output.js'))
  ;
}
```

### options

You can to pass an additional options object as a second parameter of
watchify. Its properties are:

`opts.delay` is the amount of time in milliseconds to wait before emitting
an "update" event after a change. Defaults to `100`.

`opts.ignoreWatch` ignores monitoring files for changes. If set to `true`,
then `**/node_modules/**` will be ignored. For other possible values see
Chokidar's [documentation](https://github.com/paulmillr/chokidar#path-filtering) on "ignored".

`opts.poll` enables polling to monitor for changes. If set to `true`, then
a polling interval of 100ms is used. If set to a number, then that amount of
milliseconds will be the polling interval. For more info see Chokidar's
[documentation](https://github.com/paulmillr/chokidar#performance) on
"usePolling" and "interval".
**This option is useful if you're watching an NFS volume.**

```js
var b = browserify({ cache: {}, packageCache: {} });
// watchify defaults:
b.plugin(watchify, {
  delay: 100,
  ignoreWatch: ['**/node_modules/**'],
  poll: false
});
```

## b.close()

Close all the open watch handles.

# events

## b.on('update', function (ids) {})

When the bundle changes, emit the array of bundle `ids` that changed.

## b.on('bytes', function (bytes) {})

When a bundle is generated, this event fires with the number of bytes.

## b.on('time', function (time) {})

When a bundle is generated, this event fires with the time it took to create the
bundle in milliseconds.

## b.on('log', function (msg) {})

This event fires after a bundle was created with messages of the form:

```
X bytes written (Y seconds)
```

with the number of bytes in the bundle X and the time in seconds Y.

# working with browserify transforms

If your custom transform for browserify adds new files to the bundle in a non-standard way without requiring.
You can inform Watchify about these files by emiting a 'file' event.

```
module.exports = function(file) {
  return through(
    function(buf, enc, next) {
      /*
        manipulating file content
      */
      
      this.emit("file", absolutePathToFileThatHasToBeWatched);
      
      next();
    }
  );
};
```

# install

With [npm](https://npmjs.org) do:

```
$ npm install -g watchify
```

to get the watchify command and:

```
$ npm install watchify
```

to get just the library.

# troubleshooting

## rebuilds on OS X never trigger

It may be related to a bug in `fsevents` (see [#250](https://github.com/substack/watchify/issues/205#issuecomment-98672850)
and [stackoverflow](http://stackoverflow.com/questions/26708205/webpack-watch-isnt-compiling-changed-files/28610124#28610124)).
Try the `--poll` flag
and/or renaming the project's directory - that might help.

## watchify swallows errors

To ensure errors are reported you have to add a event listener to your bundle stream. For more information see ([browserify/browserify#1487 (comment)](https://github.com/browserify/browserify/issues/1487#issuecomment-173357516) and [stackoverflow](https://stackoverflow.com/a/22389498/1423220))

**Example:**
```
var b = browserify();
b.bundle()
  .on('error', console.error)
   ...
;
```

# see also

- [budo](https://www.npmjs.com/package/budo) – a simple development server built on watchify
- [errorify](https://www.npmjs.com/package/errorify) – a plugin to add error handling to watchify development
- [watchify-request](https://www.npmjs.com/package/watchify-request) – wraps a `watchify` instance to avoid stale bundles in HTTP requests
- [watchify-middleware](https://www.npmjs.com/package/watchify-middleware) – similar to `watchify-request`, but includes some higher-level features

# license

MIT
