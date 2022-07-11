package compute

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/flowswiss/cli/v2/internal/commands"
	"github.com/flowswiss/cli/v2/pkg/api/compute"
	"github.com/flowswiss/cli/v2/pkg/filter"
)

func KeyPairCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "key-pair",
		Aliases: []string{"key-pairs"},
		Short:   "Manage compute key pairs",
	}

	commands.Add(cmd, &keyPairListCommand{}, &keyPairCreateCommand{}, &keyPairDeleteCommand{})

	return cmd
}

type keyPairListCommand struct {
	filter string
}

func (k *keyPairListCommand) Run(cmd *cobra.Command, args []string) error {
	keyPairs, err := compute.NewKeyPairService(commands.Config.Client).List(cmd.Context())
	if err != nil {
		return fmt.Errorf("fetch key pairs: %w", err)
	}

	if len(k.filter) != 0 {
		keyPairs = filter.Find(keyPairs, k.filter)
	}

	return commands.PrintStdout(keyPairs)
}

func (k *keyPairListCommand) Build() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"show", "ls", "get"},
		Short:   "List compute key pairs",
		Long:    "Lists all compute key pairs.",
		RunE:    k.Run,
	}

	cmd.Flags().StringVar(&k.filter, "filter", "", "custom term to filter the results")

	return cmd
}

type keyPairCreateCommand struct {
	name      string
	publicKey string
}

func (k *keyPairCreateCommand) Run(cmd *cobra.Command, args []string) error {
	publicKey, err := os.ReadFile(k.publicKey)
	if err != nil {
		return fmt.Errorf("read public key: %w", err)
	}

	data := compute.KeyPairCreate{
		Name:      k.name,
		PublicKey: string(publicKey),
	}

	keyPair, err := compute.NewKeyPairService(commands.Config.Client).Create(cmd.Context(), data)
	if err != nil {
		return fmt.Errorf("create key pair: %w", err)
	}

	return commands.PrintStdout(keyPair)
}

func (k *keyPairCreateCommand) Build() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "create",
		Aliases: []string{"add", "new"},
		Short:   "Create a new compute key pair",
		Long:    "Creates a new compute key pair.",
		RunE:    k.Run,
	}

	cmd.Flags().StringVar(&k.name, "name", "", "name of the key pair")
	cmd.Flags().StringVar(&k.publicKey, "public-key", "", "path to the public key of the key pair in OpenSSH format")

	_ = cmd.MarkFlagRequired("name")
	_ = cmd.MarkFlagRequired("public-key")
	_ = cmd.MarkFlagFilename("public-key")

	return cmd
}

type keyPairDeleteCommand struct {
	force bool
}

func (k *keyPairDeleteCommand) Run(cmd *cobra.Command, args []string) error {
	keyPair, err := findKeyPair(cmd.Context(), args[0])
	if err != nil {
		return err
	}

	if !k.force && !commands.ConfirmDeletion("key pair", keyPair) {
		commands.Stderr.Println("aborted.")
		return nil
	}

	err = compute.NewKeyPairService(commands.Config.Client).Delete(cmd.Context(), keyPair.ID)
	if err != nil {
		return fmt.Errorf("delete key pair: %w", err)
	}

	return nil
}

func (k *keyPairDeleteCommand) Build() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "delete KEY-PAIR",
		Aliases: []string{"del", "remove", "rm"},
		Short:   "Delete a compute key pair",
		Long:    "Deletes a compute key pair.",
		Args:    cobra.ExactArgs(1),
		RunE:    k.Run,
	}

	cmd.Flags().BoolVar(&k.force, "force", false, "force the deletion of the key pair without asking for confirmation")

	return cmd
}

func findKeyPair(ctx context.Context, term string) (compute.KeyPair, error) {
	keyPairs, err := compute.NewKeyPairService(commands.Config.Client).List(ctx)
	if err != nil {
		return compute.KeyPair{}, fmt.Errorf("fetch key pairs: %w", err)
	}

	keyPair, err := filter.FindOne(keyPairs, term)
	if err != nil {
		return compute.KeyPair{}, fmt.Errorf("find key pair: %w", err)
	}

	return keyPair, nil
}
