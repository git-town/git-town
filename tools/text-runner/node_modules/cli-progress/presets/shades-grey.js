const _colors = require('colors');

// cli-progress legacy style as of 1.x
module.exports = {
    format: _colors.grey(' {bar}') + ' {percentage}% | ETA: {eta}s | {value}/{total}',
    barCompleteChar: '\u2588',
    barIncompleteChar: '\u2591'
};