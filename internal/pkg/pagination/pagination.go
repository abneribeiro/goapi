package pagination

import (
	"net/http"
	"strconv"
)

const (
	DefaultPage    = 1
	DefaultPerPage = 10
	MaxPerPage     = 100
)

type Params struct {
	Page    int
	PerPage int
	Offset  int
}

func FromRequest(r *http.Request) Params {
	query := r.URL.Query()

	page := parseIntOrDefault(query.Get("page"), DefaultPage)
	if page < 1 {
		page = DefaultPage
	}

	perPage := parseIntOrDefault(query.Get("per_page"), DefaultPerPage)
	if perPage < 1 {
		perPage = DefaultPerPage
	}
	if perPage > MaxPerPage {
		perPage = MaxPerPage
	}

	return Params{
		Page:    page,
		PerPage: perPage,
		Offset:  (page - 1) * perPage,
	}
}

func parseIntOrDefault(value string, defaultValue int) int {
	if value == "" {
		return defaultValue
	}
	parsed, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}
	return parsed
}

func CalculateTotalPages(total int64, perPage int) int {
	if total == 0 {
		return 0
	}
	pages := int(total) / perPage
	if int(total)%perPage > 0 {
		pages++
	}
	return pages
}
