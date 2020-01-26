package ext

import (
	"fmt"
	"github.com/g4s8/gitstrap/context"
)

// Extensions for gitstrap
type Extensions map[string]Ext

func (e Extensions) proxies(ctx *context.Context) ([]applicable, error) {
	pxs := make([]applicable, 0)
	for name, ext := range e {
		p, err := proxy(name, ext, ctx)
		if err != nil {
			return nil, err
		}
		pxs = append(pxs, p)
	}
	return pxs, nil
}

func (e Extensions) Apply(ctx *context.Context) error {
	pxs, err := e.proxies(ctx)
	if err != nil {
		return err
	}
	for _, p := range pxs {
		if err := p.apply(ctx); err != nil {
			return err
		}
	}
	return nil
}

type Ext map[string]interface{}

type applicable interface {
	apply(ctx *context.Context) error
}

func proxy(name string, ext Ext, ctx *context.Context) (applicable, error) {
	switch name {
	case "readme":
		readme := new(readme)
		if err := readme.parse(ctx, ext); err != nil {
			return nil, err
		}
		return readme, nil
	case "0pdd":
		e := new(zpdd)
		if err := e.parse(ctx, ext); err != nil {
			return nil, err
		}
		return e, nil
	}
	return nil, fmt.Errorf("couln't find extension %s", name)
}
