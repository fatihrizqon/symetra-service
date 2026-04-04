package util

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
)

var reservedKeys = map[string]bool{
	"page":      true,
	"page_size": true,
	"search":    true,
	"sort":      true,
	"order":     true,
}

type QueryParams struct {
	Page     int
	PageSize int
	Search   string
	Fields   []string
	SortBy   string
	SortDir  string
	Filters  map[string][]string
}

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
