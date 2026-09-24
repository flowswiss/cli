package compute

import (
	"fmt"
	"strings"

	"github.com/flowswiss/cli/v2/internal/commands"
	"github.com/flowswiss/cli/v2/pkg/api/generic"
	"github.com/flowswiss/goclient/v2/compute"
	"github.com/flowswiss/goclient/v2/core"

	"github.com/flowswiss/cli/v2/pkg/api/common"
)

type (
	CertificateCreate = compute.CertificateCreateReq
	CertificateGet    = compute.CertificateGetReq
	CertificateList   = core.Cursor
	CertificateDelete = compute.CertificateDeleteReq
)

type GenericCertificateService = generic.CRD[
	compute.Certificate,
	Certificate,
	CertificateCreate,
	CertificateGet,
	CertificateList,
	CertificateDelete,
]

func CertificateService() GenericCertificateService {
	return generic.NewCRD[
		compute.Certificate,
		Certificate,
		CertificateCreate,
		CertificateGet,
		CertificateList,
		CertificateDelete,
	](commands.Client.Compute.Certificate, func(certificate compute.Certificate) Certificate {
		return Certificate(certificate)
	})
}

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

func (c Certificate) Values() map[string]any {
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

	return map[string]any{
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
