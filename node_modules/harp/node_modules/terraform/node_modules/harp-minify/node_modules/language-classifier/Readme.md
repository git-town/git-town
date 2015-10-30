
# language-classifier

  Programming language classifier based on harthur's Bayesian [classifier](https://github.com/harthur/classifier) module.

## Installation

    $ npm install language-classifier

## Example

```js
var lang = require('language-classifier');

lang('for link in links:');
// => "python"

lang('Foo.prototype.bar = function(){}');
// => "javascript"

lang('#include <stdio.h>');
// => "c"
```

## Training

  To train simply invoke `make memory`.

## Supported languages

  - ruby
  - python
  - javascript
  - objective-c
  - html
  - css
  - shell
  - c++
  - c

## License 

  MIT
