# humanize-string [![Build Status](https://travis-ci.org/sindresorhus/humanize-string.svg?branch=master)](https://travis-ci.org/sindresorhus/humanize-string)

> Convert a camelized/dasherized/underscored string into a humanized one
> Example: `fooBar-Baz_Faz` → `Foo bar baz faz`


## Install

```
$ npm install humanize-string
```


## Usage

```js
const humanizeString = require('humanize-string');

humanizeString('fooBar');
//=> 'Foo bar'

humanizeString('foo-bar');
//=> 'Foo bar'

humanizeString('foo_bar');
//=> 'Foo bar'
```


## Related

- [camelcase](https://github.com/sindresorhus/camelcase) - Convert a dash/dot/underscore/space separated string to camelcase


## License

MIT © [Sindre Sorhus](https://sindresorhus.com)
