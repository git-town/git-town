## require-new

{"gitdown": "badge", "name": "travis"}
{"gitdown": "badge", "name": "npm-version"}

> Modules are cached after the first time they are loaded. This means (among other things) that every call to require('foo') will get exactly the same object returned, if it would resolve to the same file.

– http://nodejs.org/api/modules.html#modules_caching

`require-new` requires a new module object.

`require-new` does not affect the state or behavior of `require` method.

`require-new` has been designed to be used for module testing.

## Usage

Load `require-new` module and use it to load a module just as you would with `require`:

```js
var requireNew = require('require-new'),
    myModule = requireNew('my-module');
```

### Example

If you have a module `rand.js`:

```js
module.exports = Math.random();
```

Then requiring this module several times will result in the same response:

```js
require('./rand.js'); // 0.697190385311842
require('./rand.js'); // 0.697190385311842
```

Modules are cached in a [`require.cache`](http://nodejs.org/api/globals.html#globals_require_cache) object when they are required. 

`require-new` deletes the key value from the `require.cache` object associated with the module you are requesting, making the module reload:

```js
requireNew('./rand.js'); // 0.2123227424453944
requireNew('./rand.js'); // 0.5403654584661126
```

It then restores the state of the `require.cache` object to ensure that it does not affect the behavior of `require`:

```js
require('./rand.js'); // 0.48205413995310664
requireNew('./rand.js'); // 0.12475096038542688
requireNew('./rand.js'); // 0.2615479789674282
require('./rand.js'); // 0.48205413995310664
```

## Alternative Design

You can export a function and call that function. This will make the module execute code multiple times.

```js
module.exports = function () {
    return Math.random();
};
```

```js
require('rand.js')(); // 0.561616780469194
require('rand.js')(); // 0.6468832329846919
```

See discussion on [Stack Overflow](http://stackoverflow.com/questions/9210542/node-js-require-cache-possible-to-invalidate/11477602) to consider the pros and cons.

## Download

Download using NPM:

```sh
npm install require-new
```