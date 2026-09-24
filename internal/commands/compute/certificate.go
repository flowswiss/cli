package compute

import (
	"context"
	"encoding/base64"
	"fmt"
	"os"

	"github.com/flowswiss/goclient/v2/core"
	"github.com/spf13/cobra"

	"github.com/flowswiss/cli/v2/internal/commands"
	"github.com/flowswiss/cli/v2/pkg/api/common"
	"github.com/flowswiss/cli/v2/pkg/api/compute"
)

func CertificateCommand(app commands.Application) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "certificate",
		Short: "Manage compute certificates",
	}

	commands.Add(app, cmd,
		&certificateListCommand{},
		&certificateCreateCommand{},
		&certificateDeleteCommand{},
	)

	return cmd
}

type certificateListCommand struct {
	filter string
}

func (c *certificateListCommand) Run(cmd *cobra.Command, args []string) error {
	items, err := compute.CertificateService().Filter(cmd.Context(), core.CursorAll, c.filter)
	if err != nil {
		return fmt.Errorf("fetch certificates: %w", err)
	}

	return commands.PrintStdout(items)
}

func (c *certificateListCommand) CompleteArg(cmd *cobra.Command, args []string, toComplete string) (
	[]string,
	cobra.ShellCompDirective,
) {
	return nil, cobra.ShellCompDirectiveNoFileComp
}

func (c *certificateListCommand) Build(app commands.Application) *cobra.Command {
	cmd := &cobra.Command{
		Use:               "list",
		Aliases:           []string{"show", "ls", "get"},
		Short:             "List certificates",
		Long:              "Lists all certificates of the current tenant.",
		ValidArgsFunction: c.CompleteArg,
		RunE:              c.Run,
	}

	cmd.Flags().StringVar(&c.filter, "filter", "", "custom term to filter the results")

	return cmd
}

type certificateCreateCommand struct {
	name        string
	location    string
	certificate string
	privateKey  string
}

func (c *certificateCreateCommand) Run(cmd *cobra.Command, args []string) error {
	location, err := common.FindLocation(cmd.Context(), commands.Client, c.location)
	if err != nil {
		return err
	}

	certificate, err := os.ReadFile(c.certificate)
	if err != nil {
		return fmt.Errorf("read certificate: %w", err)
	}

	privateKey, err := os.ReadFile(c.privateKey)
	if err != nil {
		return fmt.Errorf("read private key: %w", err)
	}

	data := compute.CertificateCreate{
		Name:        c.name,
		LocationID:  location.ID,
		Certificate: base64.StdEncoding.EncodeToString(certificate),
		PrivateKey:  base64.StdEncoding.EncodeToString(privateKey),
	}

	item, err := compute.CertificateService().Create(cmd.Context(), data)
	if err != nil {
		return fmt.Errorf("create certificate: %w", err)
	}

	return commands.PrintStdout(item)
}

func (c *certificateCreateCommand) CompleteArg(cmd *cobra.Command, args []string, toComplete string) (
	[]string,
	cobra.ShellCompDirective,
) {
	return nil, cobra.ShellCompDirectiveNoFileComp
}

func (c *certificateCreateCommand) Build(app commands.Application) *cobra.Command {
	cmd := &cobra.Command{
		Use:               "create",
		Short:             "Create a certificate",
		Long:              "Creates a new certificate",
		ValidArgsFunction: c.CompleteArg,
		RunE:              c.Run,
	}

	cmd.Flags().StringVar(&c.name, "name", "", "name of the certificate")
	cmd.Flags().StringVar(&c.location, "location", "", "location of the certificate")
	cmd.Flags().StringVar(&c.certificate, "certificate", "", "path to the certificate file")
	cmd.Flags().StringVar(&c.privateKey, "private-key", "", "path to the private key file")

	_ = cmd.MarkFlagRequired("name")
	_ = cmd.MarkFlagRequired("location")
	_ = cmd.MarkFlagRequired("certificate")
	_ = cmd.MarkFlagRequired("private-key")

	_ = cmd.MarkFlagFilename("certificate")
	_ = cmd.MarkFlagFilename("private-key")

	return cmd
}

type certificateDeleteCommand struct {
	force bool
}

func (c *certificateDeleteCommand) Run(cmd *cobra.Command, args []string) error {
	certificate, err := findCertificate(cmd.Context(), args[0])
	if err != nil {
		return err
	}

	if !c.force && !commands.ConfirmDeletion("certificate", certificate) {
		commands.Stderr.Println("aborted.")
		return nil
	}

	err = compute.CertificateService().Delete(cmd.Context(), compute.CertificateDelete{
		ID: uint(certificate.ID),
	})
	if err != nil {
		return fmt.Errorf("delete certificate: %w", err)
	}

	return nil
}

func (c *certificateDeleteCommand) CompleteArg(cmd *cobra.Command, args []string, toComplete string) (
	[]string,
	cobra.ShellCompDirective,
) {
	if len(args) == 0 {
		return completeCertificate(cmd.Context(), toComplete)
	}

	return nil, cobra.ShellCompDirectiveNoFileComp
}

func (c *certificateDeleteCommand) Build(app commands.Application) *cobra.Command {
	cmd := &cobra.Command{
		Use:               "delete CERTIFICATE",
		Short:             "Delete certificate",
		Long:              "Deletes a compute certificate.",
		Args:              cobra.ExactArgs(1),
		ValidArgsFunction: c.CompleteArg,
		RunE:              c.Run,
	}

	cmd.Flags().BoolVar(&c.force, "force", false, "force the deletion of the certificate without asking for confirmation")

	return cmd
}

func completeCertificate(ctx context.Context, term string) ([]string, cobra.ShellCompDirective) {
	certificates, err := compute.CertificateService().Filter(ctx, core.CursorAll, term)
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}

	names := make([]string, len(certificates))
	for i, item := range certificates {
		names[i] = item.Name
	}

	return names, cobra.ShellCompDirectiveNoFileComp
}

func findCertificate(ctx context.Context, term string) (compute.Certificate, error) {
	certificate, err := compute.CertificateService().FindOne(ctx, core.CursorAll, term)
	if err != nil {
		return compute.Certificate{}, fmt.Errorf("fetch certificates: %w", err)
	}

	return certificate, nil
}
