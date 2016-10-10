package main

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestSplitCommand(t *testing.T) {
	assert := assert.New(t)
	params := map[string][]string{
		" foo ": {"foo"},
		"foo bar": {"foo", "bar"},
		"foo bar baz": {"foo", "bar", "baz"},
		"foo 'bar baz'": {"foo", "bar baz"},
		"foo \"bar baz\"": {"foo", "bar baz"},
		"foo \"bar baz\"  ": {"foo", "bar baz"},
		"foo \"bar 'a' baz  \"": {"foo", "bar 'a' baz  "},
	}
	for input, expected := range params {
		actual := SplitCommand(input)
		assert.EqualValues(expected, actual)
	}
}
