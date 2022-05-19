package console

import (
	"fmt"
	"io"
	"os"

	"golang.org/x/crypto/ssh/terminal"
)

const (
	AnsiReset     = "\033[0m"
	AnsiBold      = "\u001B[1m"
	AnsiItalic    = "\u001B[3m"
	AnsiUnderline = "\u001B[4m"
)

type AnsiColor int

const (
	AnsiBlack AnsiColor = iota + 30
	AnsiRed
	AnsiGreen
	AnsiYellow
	AnsiBlue
	AnsiMagenta
	AnsiCyan
	AnsiWhite

	AnsiBackground AnsiColor = 10
	AnsiBright     AnsiColor = 60
)

type Console struct {
	Writer       io.Writer
	EnableColors bool
}

func NewConsoleOutput(writer *os.File) *Console {
	return &Console{
		Writer:       writer,
		EnableColors: terminal.IsTerminal(int(writer.Fd())),
	}
}

func (c *Console) AnsiSequence(sequence string) *Console {
	if c.EnableColors {
		_, _ = fmt.Fprintf(c.Writer, sequence)
	}
	return c
}

func (c *Console) Color(color AnsiColor) *Console {
	return c.AnsiSequence(fmt.Sprintf("\033[%dm", color))
}

func (c *Console) Printf(format string, a ...interface{}) *Console {
	_, _ = fmt.Fprintf(c.Writer, format, a...)
	return c
}

func (c *Console) Print(a ...interface{}) *Console {
	_, _ = fmt.Fprint(c.Writer, a...)
	return c
}

func (c *Console) Println(a ...interface{}) *Console {
	_, _ = fmt.Fprintln(c.Writer, a...)
	return c
}

func (c *Console) Reset() *Console {
	return c.AnsiSequence(AnsiReset)
}

func (c *Console) Bold(format string, a ...interface{}) *Console {
	return c.AnsiSequence(AnsiBold).
		Printf(format, a...).
		Reset()
}

func (c *Console) Errorf(format string, a ...interface{}) *Console {
	return c.Color(AnsiRed).
		Printf(format, a...).
		Reset()
}
