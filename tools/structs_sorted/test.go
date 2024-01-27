// This file is test data.
package main

// Test1 is an unsorted struct definition.
type Unsorted1 struct {
	field2 int // this field should not be first
	field1 int // this field should not be last
}

// Test1 is an unsorted struct definition.
type Unsorted2 struct {
	// this field should not be first
	field2 int
	// this field should not be last
	field1 int
}

type Nested struct {
	Unsorted1
	Another Unsorted2
}

// Test2 is a sorted struct definition.
type Sorted struct {
	fieldA int // this field is okay
	fieldB int // this field is also okay
}

func test() {
	_ = Unsorted1{
		field2: 2,
		field1: 1,
	}

	_ = Sorted{
		fieldA: 1,
		fieldB: 2,
	}

	_ = Nested{
		Unsorted1: Unsorted1{
			field2: 2,
			field1: 1,
		},
		Another: Unsorted2{
			field2: 2,
			field1: 1,
		},
	}
}
