package context

import "fmt"

func (ctx *Context) String() string {
	return fmt.Sprintf("path=%s opts=%s", ctx.Path, ctx.Opt)
}
