# terraform

> Terraform is the pre-processor engine for the [Harp](https://github.com/sintaxi/harp) web server. Terraform does not write or serve files. It processes and provides a layout/partial paradigm.

## Features

- pre-processors
- layouts
- partials
- metadata (via _data.json)
- LRU cache (production mode)

### Supported Pre-Processors

**HTML** – EJS, Jade, Markdown
**CSS** – LESS, Stylus, Sass (SCSS)
**JavaScript** – CoffeeScript

## Install

```
npm install terraform
```

## API


Step 1) require the library

```javascript
var terraform = require('terraform')
```

Step 2) set the root

- publicPath (String): path to public directory
- globals (Object): global variables to be available to every template

```javascript
var planet = terraform.root("path/to/public/dir", { "title": "Bitchin" })
```

Step 3) render a file

```javascript
planet.render('index.jade', { "title": "Override the global title" }, function(error, body){
  console.log(body)
})
```

## Tests

Please run the tests

```
npm install
npm test
```

## License

Copyright © 2012–2014 Chloi Inc. All rights reserved.

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the “Software”), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED “AS IS”, WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
