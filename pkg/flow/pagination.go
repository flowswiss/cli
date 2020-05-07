package flow

import (
	"net/http"
	"regexp"
	"strconv"
)

type PaginationOptions struct {
	Page     int `url:"page,omitempty"`
	PerPage  int `url:"per_page,omitempty"`
	NoFilter int `url:"no_filter,omitempty"`
}

type Links struct {
	First   string
	Last    string
	Current string
	Prev    string
	Next    string
}

type Pagination struct {
	Count      int
	Limit      int
	TotalCount int

	CurrentPage int
	TotalPages  int

	Links Links
}

func parseIntOrZero(val string) int {
	i, err := strconv.ParseInt(val, 10, 32)
	if err != nil {
		return 0
	}

	return int(i)
}

func parseLinks(header string) Links {
	res := Links{}
	regex := regexp.MustCompile("<([^>]+)>; rel=\"(\\w+)\"(?:,\\s?)?")

	links := regex.FindAllStringSubmatch(header, 5)
	for _, link := range links {
		switch link[2] {
		case "first":
			res.First = link[1]
		case "last":
			res.Last = link[1]
		case "self":
			res.Current = link[1]
		case "next":
			res.Next = link[1]
		case "prev":
			res.Prev = link[1]
		default:
			continue
		}
	}

	return res
}

func parsePagination(res *http.Response) Pagination {
	return Pagination{
		Count:       parseIntOrZero(res.Header.Get("X-Pagination-Count")),
		Limit:       parseIntOrZero(res.Header.Get("X-Pagination-Limit")),
		TotalCount:  parseIntOrZero(res.Header.Get("X-Pagination-Total-Count")),
		CurrentPage: parseIntOrZero(res.Header.Get("X-Pagination-Current-Page")),
		TotalPages:  parseIntOrZero(res.Header.Get("X-Pagination-Total-Pages")),
		Links:       parseLinks(res.Header.Get("Link")),
	}
}
