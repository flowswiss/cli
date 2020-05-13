package commands

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type CommandLineCredentialsProvider struct {
}

func (t *CommandLineCredentialsProvider) Username() string {
	return authConfig.GetString("username")
}

func (t *CommandLineCredentialsProvider) Password() string {
	return authConfig.GetString("password")
}

func (t *CommandLineCredentialsProvider) TwoFactorCode() string {
	if config.TwoFactorCode != "" {
		return config.TwoFactorCode
	}

	fmt.Printf("You have enabled two factor. Please enter your verification code: ")
	reader := bufio.NewReader(os.Stdin)
	code, err := reader.ReadString('\n')
	if err != nil {
		handleError(err)
	}

	return strings.TrimSpace(code)
}
