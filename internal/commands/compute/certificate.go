package compute

import (
	"context"
	"encoding/base64"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/flowswiss/cli/v2/internal/commands"
	"github.com/flowswiss/cli/v2/pkg/api/common"
	"github.com/flowswiss/cli/v2/pkg/api/compute"
	"github.com/flowswiss/cli/v2/pkg/filter"
)

func CertificateCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "certificate",
		Short: "Manage compute certificates",
	}

	commands.Add(cmd, &certificateListCommand{}, &certificateCreateCommand{}, &certificateDeleteCommand{})

	return cmd
}

type certificateListCommand struct {
	filter string
}

func (r *certificateListCommand) Run(cmd *cobra.Command, args []string) error {
	items, err := compute.NewCertificateService(commands.Config.Client).List(cmd.Context())
	if err != nil {
		return fmt.Errorf("fetch certificates: %w", err)
	}

	if len(r.filter) != 0 {
		items = filter.Find(items, r.filter)
	}

	return commands.PrintStdout(items)
}

func (r *certificateListCommand) Build() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"show", "ls", "get"},
		Short:   "List certificates",
		Long:    "Lists all certificates of the current tenant.",
		RunE:    r.Run,
	}

	cmd.Flags().StringVar(&r.filter, "filter", "", "custom term to filter the results")

	return cmd
}

type certificateCreateCommand struct {
	name        string
	location    string
	certificate string
	privateKey  string
}

func (r *certificateCreateCommand) Run(cmd *cobra.Command, args []string) error {
	location, err := common.FindLocation(cmd.Context(), commands.Config.Client, r.location)
	if err != nil {
		return err
	}

	certificate, err := os.ReadFile(r.certificate)
	if err != nil {
		return fmt.Errorf("read certificate: %w", err)
	}

	privateKey, err := os.ReadFile(r.privateKey)
	if err != nil {
		return fmt.Errorf("read private key: %w", err)
	}

	data := compute.CertificateCreate{
		Name:        r.name,
		LocationID:  location.ID,
		Certificate: base64.StdEncoding.EncodeToString(certificate),
		PrivateKey:  base64.StdEncoding.EncodeToString(privateKey),
	}

	item, err := compute.NewCertificateService(commands.Config.Client).Create(cmd.Context(), data)
	if err != nil {
		return fmt.Errorf("create certificate: %w", err)
	}

	return commands.PrintStdout(item)
}

func (r *certificateCreateCommand) Build() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a certificate",
		Long:  "Creates a new certificate",
		RunE:  r.Run,
	}

	cmd.Flags().StringVar(&r.name, "name", "", "name of the certificate")
	cmd.Flags().StringVar(&r.location, "location", "", "location of the certificate")
	cmd.Flags().StringVar(&r.certificate, "certificate", "", "path to the certificate file")
	cmd.Flags().StringVar(&r.privateKey, "private-key", "", "path to the private key file")

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

func (r *certificateDeleteCommand) Run(cmd *cobra.Command, args []string) error {
	certificate, err := findCertificate(cmd.Context(), args[0])
	if err != nil {
		return err
	}

	if !r.force && !commands.ConfirmDeletion("certificate", certificate) {
		commands.Stderr.Println("aborted.")
		return nil
	}

	err = compute.NewCertificateService(commands.Config.Client).Delete(cmd.Context(), certificate.ID)
	if err != nil {
		return fmt.Errorf("delete certificate: %w", err)
	}

	return nil
}

func (r *certificateDeleteCommand) Build() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete CERTIFICATE",
		Short: "Delete certificate",
		Long:  "Deletes a compute certificate.",
		Args:  cobra.ExactArgs(1),
		RunE:  r.Run,
	}

	cmd.Flags().BoolVar(&r.force, "force", false, "force the deletion of the certificate without asking for confirmation")

	return cmd
}

func findCertificate(ctx context.Context, term string) (compute.Certificate, error) {
	certificates, err := compute.NewCertificateService(commands.Config.Client).List(ctx)
	if err != nil {
		return compute.Certificate{}, fmt.Errorf("fetch certificates: %w", err)
	}

	certificate, err := filter.FindOne(certificates, term)
	if err != nil {
		return compute.Certificate{}, fmt.Errorf("find certificate: %w", err)
	}

	return certificate, nil
}
