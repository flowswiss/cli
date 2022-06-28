package commands

import (
	"regexp"
	"strings"

	"github.com/spf13/cobra"
)

type CommandBuilder interface {
	Build() *cobra.Command
}

func Add(parent *cobra.Command, builder ...CommandBuilder) {
	for _, b := range builder {
		parent.AddCommand(b.Build())
	}
}

var (
	formatIndentRegex = regexp.MustCompile("\n[ \t]*")
	formatIndent      = "  "
)

func FormatAndIndent(text string, indent int) string {
	firstIndent := formatIndentRegex.FindString(text)
	totalIndent := strings.Repeat(formatIndent, indent)

	text = strings.TrimSpace(text)
	text = strings.ReplaceAll(text, firstIndent, "\n"+totalIndent)
	text = totalIndent + text
	return text
}

func FormatHelp(help string) string {
	return FormatAndIndent(help, 0)
}

func FormatExamples(examples string) string {
	return FormatAndIndent(examples, 1)
}
