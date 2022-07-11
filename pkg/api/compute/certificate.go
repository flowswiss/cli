package compute

import (
	"context"
	"fmt"
	"strings"

	"github.com/flowswiss/goclient"
	"github.com/flowswiss/goclient/compute"

	"github.com/flowswiss/cli/v2/pkg/api/common"
)

type Certificate compute.Certificate

func (c Certificate) String() string {
	return c.Name
}

func (c Certificate) Keys() []string {
	return []string{fmt.Sprint(c.ID), c.Name, c.Details.Serial}
}

func (c Certificate) Columns() []string {
	return []string{"id", "name", "location", "valid from", "valid to", "serial", "subject", "issuer"}
}

func (c Certificate) Values() map[string]interface{} {
	subjectBuffer := strings.Builder{}
	for key, value := range c.Details.Subject {
		if subjectBuffer.Len() != 0 {
			subjectBuffer.WriteString(", ")
		}

		subjectBuffer.WriteString(fmt.Sprintf("%s=%s", key, value))
	}

	issuerBuffer := strings.Builder{}
	for key, value := range c.Details.Issuer {
		if issuerBuffer.Len() != 0 {
			issuerBuffer.WriteString(", ")
		}

		issuerBuffer.WriteString(fmt.Sprintf("%s=%s", key, value))
	}

	return map[string]interface{}{
		"id":         c.ID,
		"name":       c.Name,
		"location":   common.Location(c.Location),
		"valid from": c.Details.ValidFrom,
		"valid to":   c.Details.ValidTo,
		"serial":     c.Details.Serial,
		"subject":    subjectBuffer.String(),
		"issuer":     issuerBuffer.String(),
	}
}

type CertificateService struct {
	delegate compute.CertificateService
}

func NewCertificateService(client goclient.Client) CertificateService {
	return CertificateService{
		delegate: compute.NewCertificateService(client),
	}
}

func (c CertificateService) List(ctx context.Context) ([]Certificate, error) {
	res, err := c.delegate.List(ctx, goclient.Cursor{NoFilter: 1})
	if err != nil {
		return nil, err
	}

	items := make([]Certificate, len(res.Items))
	for idx, item := range res.Items {
		items[idx] = Certificate(item)
	}

	return items, nil
}

type CertificateCreate = compute.CertificateCreate

func (c CertificateService) Create(ctx context.Context, data CertificateCreate) (Certificate, error) {
	res, err := c.delegate.Create(ctx, data)
	if err != nil {
		return Certificate{}, err
	}

	return Certificate(res), nil
}

func (c CertificateService) Delete(ctx context.Context, id int) error {
	return c.delegate.Delete(ctx, id)
}
