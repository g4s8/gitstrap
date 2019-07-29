package gitstrap

import (
	"context"
	"os"
	"testing"

	"github.com/google/go-github/github"
	"github.com/stretchr/testify/assert"
	"gopkg.in/h2non/gock.v1"
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

func Test_getOwner(t *testing.T) {
	assert := assert.New(t)

	ctx := &strapCtx{
		cli: github.NewClient(nil),
		ctx: context.Background(),
	}
	opt := make(Options)

	// Failed request
	gock.DisableNetworking()
	owner, err := getOwner(ctx, opt)
	assert.Equal("", owner)
	assert.NotNil(err)
	// Successful request
	// gock.New("https://api.github.com").Get("/user").
	// 	Reply(200).JSON(map[string]string{"login": "tenderlove"})
	// defer gock.Off()

	// owner, err = getOwner(ctx, opt)
	// assert.Equal("tenderlove", owner)
	// assert.Nil(err)

	// // Organization
	// opt["org"] = "nasa"
	// owner, err = getOwner(ctx, opt)
	// assert.Equal("nasa", owner)
	// assert.Nil(err)
}
