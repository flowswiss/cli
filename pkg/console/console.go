package console

import (
	"fmt"
	"io"
	"os"

	"golang.org/x/term"
)

const (
	Reset     = "\033[0m"
	Bold      = "\u001B[1m"
	Italic    = "\u001B[3m"
	Underline = "\u001B[4m"
)

type Color int

const (
	Black Color = iota + 30
	Red
	Green
	Yellow
	Blue
	Magenta
	Cyan
	White

	Background Color = 10
	Bright     Color = 60
)

type Writer interface {
	io.Writer

	Color(color Color) Writer
	Bold() Writer
	Reset() Writer

	Printf(format string, a ...interface{}) Writer
	Print(a ...interface{}) Writer
	Println(a ...interface{}) Writer

	Errorf(format string, a ...interface{}) Writer
}

func NewConsoleOutput(file *os.File) Writer {
	if term.IsTerminal(int(file.Fd())) {
		return ansiWriter{File: file}
	}

	return plainWriter{File: file}
}

type plainWriter struct {
	*os.File
}

func (w plainWriter) Color(Color) Writer { return w }
func (w plainWriter) Bold() Writer       { return w }
func (w plainWriter) Reset() Writer      { return w }

func (w plainWriter) Printf(format string, a ...interface{}) Writer {
	_, _ = fmt.Fprintf(w.File, format, a...)
	return w
}

func (w plainWriter) Print(a ...interface{}) Writer {
	_, _ = fmt.Fprint(w.File, a...)
	return w
}

func (w plainWriter) Println(a ...interface{}) Writer {
	_, _ = fmt.Fprintln(w.File, a...)
	return w
}

func (w plainWriter) Errorf(format string, a ...interface{}) Writer { return w.Printf(format, a...) }

type ansiWriter struct {
	*os.File
}

func (w ansiWriter) Color(color Color) Writer {
	return w.Printf("\033[%dm", color)
}

func (w ansiWriter) Bold() Writer {
	return w.Print(Bold)
}

func (w ansiWriter) Reset() Writer {
	return w.Print(Reset)
}

func (w ansiWriter) Printf(format string, a ...interface{}) Writer {
	_, _ = fmt.Fprintf(w.File, format, a...)
	return w
}

func (w ansiWriter) Print(a ...interface{}) Writer {
	_, _ = fmt.Fprint(w.File, a...)
	return w
}

func (w ansiWriter) Println(a ...interface{}) Writer {
	_, _ = fmt.Fprintln(w.File, a...)
	return w
}

func (w ansiWriter) Errorf(format string, a ...interface{}) Writer {
	return w.Color(Red).
		Printf(format, a...).
		Reset()
}
