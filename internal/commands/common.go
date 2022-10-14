package commands

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/spf13/cobra"

	"github.com/flowswiss/cli/v2/pkg/api/common"
	"github.com/flowswiss/cli/v2/pkg/console"
)

type CommandBuilder interface {
	Build(app Application) *cobra.Command
}

func Add(app Application, parent *cobra.Command, builder ...CommandBuilder) {
	for _, b := range builder {
		parent.AddCommand(b.Build(app))
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

func Confirm(message string) bool {
	return console.Confirm(Stderr, message)
}

func ConfirmDeletion(kind string, item fmt.Stringer) bool {
	return console.Confirm(Stderr, fmt.Sprintf("Are you sure you want to delete the %s %q?", kind, item))
}

func WaitForOrder(ctx context.Context, action string, ordering common.Ordering) (common.Order, error) {
	progress := console.NewProgress(action)
	defer progress.Done()

	go progress.Display(Stderr)

	return common.WaitForOrder(ctx, Config.Client, ordering)
}
