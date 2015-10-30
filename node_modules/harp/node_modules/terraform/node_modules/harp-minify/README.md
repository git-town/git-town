# minify

Super clean API for minifying Javascript, HTML or CSS.

So you don't have to keep googling for the right tool or the tool’s API. And so that you get a nice CLI regardless.

This is [Harp](https://github.com/sintaxi/harp)’s fork, which is nearly identical to [the original](https://github.com/ianstormtaylor/minify). The differences are:

- It’s published to npm, to remove the dependency on GitHub for installing
- It has more up-to-date dependencies
- It has npm run scripts for development

## Installation

```sh
npm install harp-minify
```

## CLI

```
Usage: minify [<input>] [<output>]

Options:

  -h, --help     output usage information
  -V, --version  output the version number

Examples:

  # pass an input and output file
  $ minify input.css output.css

  # use stdin and stdout
  $ cat input.css | myth | minify > output.css
```

## Node

```javascript
var minify = require('minify');

// choose javascript, html or css
var js = minify.js('js string');
var html = minify.html('html string');
var css = minify.css('css string');

// or pass an unknown string
var min = minify('unknown string');
```

When using JavaScript, you may also alter the default options using the same API as [UglifyJS](https://github.com/mishoo/UglifyJS2):


```javascript
var js = minify.js('js string', {
  compress: false,
  mangle: false
});
```

## License

The MIT License (MIT)

Copyright © 2013–2015, Ian Storm Taylor &lt;ian@ianstormtaylor.com&gt;
Copyright © 2015 [Chloi Inc.](http://chloi.io)


Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the “Software”), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED “AS IS,” WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
