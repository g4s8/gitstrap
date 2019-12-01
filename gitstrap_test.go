package gitstrap

import (
	"context"
	"github.com/g4s8/gopwd"
	"github.com/google/go-github/github"
	"github.com/stretchr/testify/assert"
	"gopkg.in/h2non/gock.v1"
	"os"
	"testing"
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
	cfg.expand()
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

// This test will confirm that gitstrap can get the use the Work Directory name if the reponame not given in the config
func Test_getRepoNameIfNotInConfig(t *testing.T) {
	name, err := gopwd.Name()

	if "gitstrap" != name {
		t.Errorf("Failed to get repo name , expacted : gitstrap got:  %s , %s", name, err)
	}
}
