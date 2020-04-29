'use strict';
const decamelize = require('decamelize');

const humanizeString = input => {
	if (typeof input !== 'string') {
		throw new TypeError('Expected a string');
	}

	input = decamelize(input);
	input = input.toLowerCase().replace(/[_-]+/g, ' ').replace(/\s{2,}/g, ' ').trim();
	input = input.charAt(0).toUpperCase() + input.slice(1);

	return input;
};

module.exports = humanizeString;
// TODO: Remove this for the next major release
module.exports.default = humanizeString;
