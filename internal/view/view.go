package view

import (
	"fmt"
	"io"
	"os"
)

// Printer - lines printer
type Printer interface {
	Print(string)
}

// Console printer
var Console Printer

type stdOutPrinter struct {
	w io.Writer
}

func (p *stdOutPrinter) Print(line string) {
	fmt.Fprintln(p.w, line)
}

func init() {
	Console = &stdOutPrinter{os.Stdout}
}

// Printable - entry that can be printed on printer
type Printable interface {
	PrintOn(Printer)
}

func RenderOn(v Printer, items <-chan Printable, errs <-chan error) error {
	for {
		select {
		case next, ok := <-items:
			if ok {
				next.PrintOn(v)
			} else {
				return nil
			}
		case err, ok := <-errs:
			if ok {
				return err
			} else {
				return nil
			}
		}
	}
}
