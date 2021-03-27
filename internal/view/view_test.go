package view

import (
	"errors"
	"fmt"
	"testing"

	m "github.com/g4s8/go-matchers"
)

type item struct {
	x int
}

func (i *item) PrintOn(p Printer) {
	p.Print(fmt.Sprintf("item-%d", i.x))
}

type testPrinter struct {
	msgs []string
}

func (p *testPrinter) Print(msg string) {
	p.msgs = append(p.msgs, msg)
}

func Test_RenderOnSuccess(t *testing.T) {
	flow := make(chan Printable)
	errs := make(chan error)
	done := make(chan struct{})
	target := new(testPrinter)
	go func() {
		_ = RenderOn(target, flow, errs)
		close(done)
	}()
	go func() {
		for i := 0; i < 5; i++ {
			flow <- &item{i}
		}
		close(flow)
		close(errs)
	}()
	<-done
	assert := m.Assert(t)
	assert.That("Printer received messages in order",
		target.msgs,
		m.Eq([]string{"item-0", "item-1", "item-2", "item-3", "item-4"}))
}

func Test_renderOnError(t *testing.T) {
	flow := make(chan Printable)
	errs := make(chan error)
	done := make(chan struct{})
	target := new(testPrinter)
	printerErrs := make(chan error, 1)
	expect := errors.New("printing-error")
	go func() {
		if err := RenderOn(target, flow, errs); err != nil {
			fmt.Println("renderOn received error")
			printerErrs <- err
		}
		close(done)
	}()
	go func() {
		errs <- expect
		close(flow)
		close(errs)
	}()
	fmt.Println("wait")
	<-done
	fmt.Println("done")
	var actual error
	select {
	case e := <-printerErrs:
		actual = e
	default:
		actual = nil
	}
	fmt.Println("err")
	assert := m.Assert(t)
	assert.That("RenderOn fails with error", actual, m.AllOf(m.Not(m.Nil()), m.Eq(expect)))
}
