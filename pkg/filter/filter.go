package filter

import (
	"fmt"
	"strings"
)

type Filterable interface {
	Keys() []string
}

func Find[T Filterable](items []T, term string) []T {
	term = strings.ToLower(term)

	var filtered []T
	for _, item := range items {
		identifiers := item.Keys()
		for _, identifier := range identifiers {
			if strings.Contains(strings.ToLower(identifier), term) {
				filtered = append(filtered, item)
				break
			}
		}
	}

	return filtered
}

func FindOne[T Filterable](items []T, term string) (res T, err error) {
	var filtered = Find[T](items, term)

	if len(filtered) == 0 {
		return res, fmt.Errorf("no item found searching for the term %q", term)
	}

	if len(filtered) > 1 {
		return res, fmt.Errorf("term %q resulted in multiple results", term)
	}

	return filtered[0], nil
}
