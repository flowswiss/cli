package console

import (
	"bufio"
	"os"
	"strings"

	"golang.org/x/term"
)

const (
	confirmNo confirmOption = iota
	confirmYes
)

var (
	confirmOptionNames = map[confirmOption]string{
		confirmNo:  "n",
		confirmYes: "y",
	}
)

type confirmOption int

func (c confirmOption) String() string {
	return confirmOptionNames[c]
}

type optConstraint interface {
	comparable
	String() string
}

var reader = bufio.NewReader(os.Stdin)

func Confirm(writer Writer, question string) bool {
	res, err := Ask(writer, question, confirmNo, confirmYes)
	if err != nil {
		return false
	}

	return res == confirmYes
}

func Ask[T optConstraint](writer Writer, question string, opts ...T) (res T, err error) {
	if len(opts) == 0 {
		panic("no selection provided")
	}

	var defaultOpt T
	ok := false

	for !ok {
		writer.Print(question)

		writer.Print(" [")
		lookup := map[string]T{}

		for i, opt := range opts {
			if i != 0 {
				writer.Print("/")
			}

			name := opt.String()
			lookup[name] = opt

			if opt == defaultOpt {
				name = strings.ToUpper(name)
			}

			writer.Printf("%s", name)
		}
		writer.Print("] ")

		answer, err := reader.ReadString('\n')
		if err != nil {
			return res, err
		}

		answer = strings.TrimSpace(answer)
		res, ok = lookup[answer]

		if answer == "" {
			res = defaultOpt
			break
		}
	}

	return res, nil
}

func Password(writer Writer, prompt string, valid func(string) error) (string, error) {
	for {
		writer.Print(prompt)
		writer.Print(": ")

		password, err := term.ReadPassword(int(os.Stdin.Fd()))
		if err != nil {
			return "", err
		}

		writer.Println()

		err = valid(string(password))
		if err == nil {
			return string(password), nil
		}

		writer.Errorf("%v\n", err)
	}
}
