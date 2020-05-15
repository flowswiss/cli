package commands

import (
	"context"
	"fmt"
	"github.com/flowswiss/cli/internal/commands/dto"
	"github.com/flowswiss/cli/pkg/flow"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
)

var (
	keyPairCommand = &cobra.Command{
		Use:   "key-pair",
		Short: "Manage your ssh key pairs",
	}

	keyPairListCommand = &cobra.Command{
		Use:   "list",
		Short: "List all key pairs",
		RunE:  listKeyPair,
	}

	keyPairUploadCommand = &cobra.Command{
		Use:   "upload <file>",
		Short: "Upload key pair",
		Args:  cobra.ExactValidArgs(1),
		RunE:  uploadKeyPair,
	}
)

func init() {
	keyPairCommand.AddCommand(keyPairListCommand)
	keyPairCommand.AddCommand(keyPairUploadCommand)

	keyPairUploadCommand.Flags().StringP(flagName, "n", "", "custom name for your key pair")
}

func findKeyPair(filter string) (*flow.KeyPair, error) {
	keyPairs, _, err := client.KeyPair.List(context.Background(), flow.PaginationOptions{NoFilter: 1})
	if err != nil {
		return nil, err
	}

	keyPair, err := findOne(keyPairs, filter, 2)
	if err != nil {
		return nil, fmt.Errorf("key-pair: %v", err)
	}

	return keyPair.(*flow.KeyPair), nil
}

func findKeyPairByFingerprint(fingerprint string) (*flow.KeyPair, error) {
	keyPairs, _, err := client.KeyPair.List(context.Background(), flow.PaginationOptions{NoFilter: 1})
	if err != nil {
		return nil, err
	}

	for _, keyPair := range keyPairs {
		if keyPair.Fingerprint == fingerprint {
			return keyPair, nil
		}
	}

	return nil, nil
}

func listKeyPair(cmd *cobra.Command, args []string) error {
	keyPairs, _, err := client.KeyPair.List(context.Background(), flow.PaginationOptions{NoFilter: 1})
	if err != nil {
		return err
	}

	var displayable []*dto.KeyPair
	for _, keyPair := range keyPairs {
		displayable = append(displayable, &dto.KeyPair{KeyPair: keyPair})
	}

	return display(displayable)
}

func uploadKeyPair(cmd *cobra.Command, args []string) error {
	data, err := ioutil.ReadFile(args[0])
	if err != nil {
		return err
	}

	publicKey, comment, _, _, err := ssh.ParseAuthorizedKey(data)
	if err != nil {
		return err
	}

	if publicKey.Type() != ssh.KeyAlgoRSA {
		return fmt.Errorf("currently only rsa key formats are supported")
	}

	fingerprint := ssh.FingerprintLegacyMD5(publicKey)

	keyPair, err := findKeyPairByFingerprint(fingerprint)
	if err != nil {
		return err
	}

	if keyPair != nil {
		return fmt.Errorf("a key pair with this fingerprint has already been uploaded as %q", keyPair.Name)
	}

	name, err := cmd.Flags().GetString(flagName)
	if err != nil {
		return err
	}

	if name == "" {
		name = comment
	}

	if name == "" {
		return fmt.Errorf("please provide a name for your key pair")
	}

	keyPair, _, err = client.KeyPair.Create(context.Background(), &flow.KeyPairCreate{
		Name:      name,
		PublicKey: string(data),
	})

	if err != nil {
		return err
	}

	return display(&dto.KeyPair{KeyPair: keyPair})
}
