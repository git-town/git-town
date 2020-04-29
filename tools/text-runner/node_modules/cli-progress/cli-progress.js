const _SingleBar = require('./lib/single-bar');
const _MultiBar = require('./lib/multi-bar');
const _Presets = require('./presets/index');

// sub-module access
module.exports = {
    Bar: _SingleBar,
    SingleBar: _SingleBar,
    MultiBar: _MultiBar,
    Presets: _Presets
};