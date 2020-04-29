'use strict';

module.exports = function( array, value, count ){
	if (Array.isArray(this)) {
		count = value;
		value = array;
		array = this;
	}

	var index;
	var i = 0;

	while ((!count || i++ < count) && ~(index = array.indexOf(value)))
		array.splice(index, 1);

	return array;
};