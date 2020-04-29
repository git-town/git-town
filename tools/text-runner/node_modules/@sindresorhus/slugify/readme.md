# slugify [![Build Status](https://travis-ci.org/sindresorhus/slugify.svg?branch=master)](https://travis-ci.org/sindresorhus/slugify)

> Slugify a string

Useful for URLs, filenames, and IDs.

It handles most major languages, including [German (umlauts)](https://en.wikipedia.org/wiki/Germanic_umlaut), Vietnamese, Arabic, Russian, [and more](https://github.com/sindresorhus/transliterate#supported-languages).

## Install

```
$ npm install @sindresorhus/slugify
```

## Usage

```js
const slugify = require('@sindresorhus/slugify');

slugify('I ♥ Dogs');
//=> 'i-love-dogs'

slugify('  Déjà Vu!  ');
//=> 'deja-vu'

slugify('fooBar 123 $#%');
//=> 'foo-bar-123'

slugify('я люблю единорогов');
//=> 'ya-lyublyu-edinorogov'
```

## API

### slugify(string, options?)

#### string

Type: `string`

String to slugify.

#### options

Type: `object`

##### separator

Type: `string`\
Default: `'-'`

```js
const slugify = require('@sindresorhus/slugify');

slugify('BAR and baz');
//=> 'bar-and-baz'

slugify('BAR and baz', {separator: '_'});
//=> 'bar_and_baz'
```

##### lowercase

Type: `boolean`\
Default: `true`

Make the slug lowercase.

```js
const slugify = require('@sindresorhus/slugify');

slugify('Déjà Vu!');
//=> 'deja-vu'

slugify('Déjà Vu!', {lowercase: false});
//=> 'Deja-Vu'
```

##### decamelize

Type: `boolean`\
Default: `true`

Convert camelcase to separate words. Internally it does `fooBar` → `foo bar`.

```js
const slugify = require('@sindresorhus/slugify');

slugify('fooBar');
//=> 'foo-bar'

slugify('fooBar', {decamelize: false});
//=> 'foobar'
```

##### customReplacements

Type: `Array<string[]>`\
Default: `[
	['&', ' and '],
	['🦄', ' unicorn '],
	['♥', ' love ']
]`

Add your own custom replacements.

The replacements are run on the original string before any other transformations.

This only overrides a default replacement if you set an item with the same key, like `&`.

```js
const slugify = require('@sindresorhus/slugify');

slugify('Foo@unicorn', {
	customReplacements: [
		['@', 'at']
	]
});
//=> 'fooatunicorn'
```

Add a leading and trailing space to the replacement to have it separated by dashes:

```js
const slugify = require('@sindresorhus/slugify');

slugify('foo@unicorn', {
	customReplacements: [
		['@', ' at ']
	]
});
//=> 'foo-at-unicorn'
```

Another example:

```js
const slugify = require('@sindresorhus/slugify');

slugify('I love 🐶', {
	customReplacements: [
		['🐶', 'dogs']
	]
});
//=> 'i-love-dogs'
```

##### preserveLeadingUnderscore

Type: `boolean`\
Default: `false`

If your string starts with an underscore, it will be preserved in the slugified string.

Sometimes leading underscores are intentional, for example, filenames representing hidden paths on a website.

```js
const slugify = require('@sindresorhus/slugify');

slugify('_foo_bar');
//=> 'foo-bar'

slugify('_foo_bar', {preserveLeadingUnderscore: true});
//=> '_foo-bar'
```

## Related

- [slugify-cli](https://github.com/sindresorhus/slugify-cli) - CLI for this module
- [transliterate](https://github.com/sindresorhus/transliterate) - Convert Unicode characters to Latin characters using transliteration
- [filenamify](https://github.com/sindresorhus/filenamify) - Convert a string to a valid safe filename
