package util

import (
	"fmt"

	"github.com/fatihrizqon/symetra-service/internal/delivery/http/response"
)

func GenerateMeta(baseURL, search string, page, pageSize, totalCount int, filters map[string]string) response.Meta {
	if page < 1 {
		page = 1
	}

	totalPages := (totalCount + pageSize - 1) / pageSize

	buildURL := func(page int) string {
		query := fmt.Sprintf("?page=%d&page_size=%d", page, pageSize)
		if search != "" {
			query += "&search=" + search
		}
		for key, value := range filters {
			query += fmt.Sprintf("&%s=%s", key, value)
		}
		return baseURL + query
	}

	currentPage := buildURL(page)
	prevPage := getPageURL(page-1, buildURL)
	nextPage := getPageURL(page+1, buildURL, (page*pageSize) < totalCount)

	startIndex, endIndex := calculateIndices(page, pageSize, totalCount)

	return response.Meta{
		Search:     search,
		Info:       fmt.Sprintf("Showing %d to %d from %d item(s).", startIndex, endIndex, totalCount),
		Page:       page,
		TotalCount: totalCount,
		TotalPages: totalPages,
		PageSize:   pageSize,
		Links: response.Links{
			CurrentPage: currentPage,
			FirstPage:   buildURL(1),
			LastPage:    buildURL(totalPages),
			NextPage:    nextPage,
			PrevPage:    prevPage,
		},
	}
}

func getPageURL(page int, buildURL func(int) string, condition ...bool) *string {
	if len(condition) == 0 || condition[0] {
		url := buildURL(page)
		return &url
	}
	return nil
}

func calculateIndices(page, pageSize, totalCount int) (int, int) {
	startIndex := (page-1)*pageSize + 1
	endIndex := min(startIndex+pageSize-1, totalCount)
	return startIndex, endIndex
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
