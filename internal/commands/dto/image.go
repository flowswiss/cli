package dto

import (
	"bytes"
	"fmt"
	"github.com/flowswiss/cli/pkg/flow"
	"regexp"
	"strings"
)

type Image struct {
	*flow.Image
	Locations []*flow.Location
}

func (i *Image) Columns() []string {
	return []string{"id", "operating system", "version", "type", "availability", "license"}
}

func (i *Image) Values() map[string]interface{} {
	licenseBuf := &bytes.Buffer{}
	for idx, license := range i.RequiredLicenses {
		licenseBuf.WriteString(fmt.Sprintf("%s (%s)", license.Name, getPricePerHour(license.Price)))

		if idx+1 < len(i.RequiredLicenses) {
			licenseBuf.WriteString(", ")
		}
	}

	availabilityBuf := &bytes.Buffer{}
	for idx, availability := range i.AvailableLocations {
		found := false
		for _, location := range i.Locations {
			if location.Id == availability {
				availabilityBuf.WriteString(location.Name)
				found = true
				break
			}
		}

		if found && idx+1 < len(i.AvailableLocations) {
			licenseBuf.WriteString(", ")
		}
	}

	t := regexp.MustCompile("(?:_|^)[a-z]").ReplaceAllStringFunc(i.Type, func(s string) string {
		if s[0] == '_' {
			s = " " + s[1:]
		}

		return strings.ToUpper(s)
	})

	return map[string]interface{}{
		"id":               i.Id,
		"operating system": i.OperatingSystem,
		"version":          i.Version,
		"type":             t,
		"availability":     availabilityBuf.String(),
		"license":          licenseBuf.String(),
	}
}
