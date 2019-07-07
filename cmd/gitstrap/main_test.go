package main

import (
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_getToken(t *testing.T) {
	assert := assert.New(t)

	file, err := ioutil.TempFile("", "token")
	assert.Nil(err)
	_, err = file.WriteString("token from file")
	assert.Nil(err)

	token, err := getToken("token from flag", file.Name())
	assert.Equal("token from flag", token)
	assert.Nil(err)

	token, err = getToken("", file.Name())
	assert.Equal("token from file", token)
	assert.Nil(err)

	token, err = getToken("", "nonexistent.txt")
	assert.Equal("", token)
	assert.NotNil(err)
}
