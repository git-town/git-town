declare const humanizeString: {
	/**
	Convert a camelized/dasherized/underscored string into a humanized one: `fooBar-Baz_Faz` → `Foo bar baz faz`.

	@param text - The string to make human readable.

	@example
	```
	import humanizeString = require('humanize-string');

	humanizeString('fooBar');
	//=> 'Foo bar'

	humanizeString('foo-bar');
	//=> 'Foo bar'

	humanizeString('foo_bar');
	//=> 'Foo bar'
	```
	*/
	(text: string): string;

	// TODO: Remove this for the next major release, refactor the whole definition to:
	// declare function humanizeString(text: string): string;
	// export = humanizeString;
	default: typeof humanizeString;
};

export = humanizeString;
