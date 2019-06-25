package gitstrap

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_renderEnvVars(t *testing.T) {
	assert := assert.New(t)

	assert.Equal("---true---", renderEnvVars("---$CI---"))
}
