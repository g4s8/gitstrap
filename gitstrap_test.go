package gitstrap

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExpandTemplates(t *testing.T) {
	cfg := new(Config)
	if err := cfg.ParseFile(".gitstrap.yaml"); err != nil {
		t.Errorf("Failed to read config: %s", err)
	}
	cfg.Gitstrap.Templates[0].Location = "${FOO}/$BAR/baz"
	assert := assert.New(t)
	os.Setenv("FOO", "one")
	os.Setenv("BAR", "two")
	cfg.Expand()
	assert.Equal("one/two/baz", cfg.Gitstrap.Templates[0].Location)
}
