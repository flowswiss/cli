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
		if matches(item, term) != noMatch {
			filtered = append(filtered, item)
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
		var bestMatch T
		ambiguous := true

		// find the best match (the one that matches exactly the term)
		for _, item := range filtered {
			if matches(item, term) == exactMatch {
				if ambiguous {
					bestMatch = item
					ambiguous = false
				} else {
					// multiple matches found, return an error
					return res, ambiguousError(term)
				}
			}
		}

		if ambiguous {
			return res, ambiguousError(term)
		}

		return bestMatch, nil
	}

	return filtered[0], nil
}

type ambiguousError string

func (a ambiguousError) Error() string {
	return fmt.Sprintf("term %q is ambiguous", string(a))
}

type matchType int

const (
	noMatch matchType = iota
	exactMatch
	partialMatch
)

func matches[T Filterable](item T, term string) matchType {
	identifiers := item.Keys()
	for _, identifier := range identifiers {
		identifier = strings.ToLower(identifier)

		if strings.EqualFold(identifier, term) {
			return exactMatch
		}

		if strings.Contains(identifier, term) {
			return partialMatch
		}
	}

	return noMatch
}
