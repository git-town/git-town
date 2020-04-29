'use strict';

var removeValue = require('./');
var test = require('tap').test;

test('remove-value', function( t ) {
	t.deepEqual(removeValue([ 'apple', 'lemon', 'banana', 'lemon' ], 'lemon'), [ 'apple', 'banana' ] );
	t.deepEqual(removeValue([ 'apple', 'lemon', 'banana', 'lemon' ], 'lemon', 1), [ 'apple', 'banana', 'lemon' ] );
	t.deepEqual(removeValue([ 'apple', 'lemon' ], 'not there', 1), [ 'apple', 'lemon' ] );

	Array.prototype.remove = removeValue;

	var list = [ 'apple', 'lemon', 'banana' ];

	t.deepEqual(list.remove('banana'), [ 'apple', 'lemon' ]);

	t.end();
});