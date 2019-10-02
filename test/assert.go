package test

import (
	"io/ioutil"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

func assertFileExists(t *testing.T, dir, filename string) {
	assertFileExistsWithContent(t, dir, filename, filename+" content")
}

func assertFileExistsWithContent(t *testing.T, dir, filename, expectedContent string) {
	fileContent, err := ioutil.ReadFile(path.Join(dir, filename))
	assert.Nil(t, err)
	assert.Equal(t, expectedContent, string(fileContent))
}
