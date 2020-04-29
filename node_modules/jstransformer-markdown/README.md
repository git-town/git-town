# jstransformer-markdown

[Markdown](http://npm.im/markdown) support for [JSTransformers](http://github.com/jstransformers).

[![Build Status](https://img.shields.io/travis/jstransformers/jstransformer-markdown/master.svg)](https://travis-ci.org/jstransformers/jstransformer-markdown)
[![Coverage Status](https://img.shields.io/codecov/c/github/jstransformers/jstransformer-markdown/master.svg)](https://codecov.io/gh/jstransformers/jstransformer-markdown)
[![Dependency Status](https://img.shields.io/david/jstransformers/jstransformer-markdown/master.svg)](http://david-dm.org/jstransformers/jstransformer-markdown)
[![Greenkeeper badge](https://badges.greenkeeper.io/jstransformers/jstransformer-markdown.svg)](https://greenkeeper.io/)
[![NPM version](https://img.shields.io/npm/v/jstransformer-markdown.svg)](https://www.npmjs.org/package/jstransformer-markdown)

## Installation

    npm install jstransformer-markdown

## API

```js
var markdown = require('jstransformer')(require('jstransformer-markdown'))

markdown.render('# Hello World!').body
//=> '<h1>Hello World!</h1>'
```

## License

MIT
