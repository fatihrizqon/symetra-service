package util

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
)

// Reserved query parameter keys that are handled explicitly.
// Any other key is treated as a filter and placed into Filters.
var reservedKeys = map[string]bool{
	"page":      true,
	"page_size": true,
	"search":    true,
	"sort":      true,
	"order":     true,
}

// QueryParams is the single contract passed across handler → service → repository
// for all list/pagination queries. Add fields to this struct instead of changing
// function signatures when new query capabilities are needed.
type QueryParams struct {
	// Pagination
	Page     int
	PageSize int

	// Full-text search value and the DB columns to search against.
	// Fields is populated by the handler from entity.SearchableFields().
	Search string
	Fields []string

	// Sorting — validated against a per-repository whitelist to prevent injection.
	SortBy  string // raw column name requested by the caller
	SortDir string // "asc" or "desc"

	// Filters holds every query parameter that is NOT one of the reserved keys
	// above. Repository layers pass these to entity.ApplyFilters().
	// Multi-value filters are supported via repeated params: ?code=1&code=2
	Filters map[string][]string
}

// ParseQueryParams builds a QueryParams from a Fiber request context.
// fields should come from entity.SearchableFields() at the call site.
func ParseQueryParams(ctx *fiber.Ctx, fields []string) *QueryParams {
	qp := &QueryParams{
		Page:     1,
		PageSize: 10,
		Fields:   fields,
		Filters:  make(map[string][]string),
	}

	if p := ctx.Query("page"); p != "" {
		if v, err := strconv.Atoi(p); err == nil && v > 0 {
			qp.Page = v
		}
	}

	if s := ctx.Query("page_size"); s != "" {
		if v, err := strconv.Atoi(s); err == nil && v > 0 {
			qp.PageSize = v
		}
	}

	qp.Search = ctx.Query("search")

	if sort := ctx.Query("sort"); sort != "" {
		qp.SortBy = sort
	}

	if ctx.Query("order") == "desc" {
		qp.SortDir = "desc"
	} else {
		qp.SortDir = "asc"
	}

	// FIXED: handle multi-value filters
	args := ctx.Context().QueryArgs()
	args.VisitAll(func(key, value []byte) {
		k := string(key)
		if !reservedKeys[k] {
			qp.Filters[k] = append(qp.Filters[k], string(value))
		}
	})

	return qp
}
