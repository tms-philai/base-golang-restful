package utils

import (
	"base-golang-restful-app/models"
	"fmt"
	"math"
)

type PaginationHelper struct {
	baseURL string
}

func NewPaginationHelper(baseURL string) *PaginationHelper {
	return &PaginationHelper{baseURL: baseURL}
}

func (h *PaginationHelper) CreatePaginationResponse(
	data interface{},
	page int,
	pageSize int,
	totalCount int64,
	withLinks bool,
) *models.PaginationResponse {
	totalPages := int(math.Ceil(float64(totalCount) / float64(pageSize)))

	meta := models.PaginationMeta{
		Page:       page,
		PageSize:   pageSize,
		TotalCount: totalCount,
		TotalPages: totalPages,
		HasNext:    page < totalPages,
		HasPrev:    page > 1,
	}

	response := &models.PaginationResponse{
		Data:       data,
		Pagination: meta,
	}

	if withLinks && h.baseURL != "" {
		response.Links = h.generateLinks(page, totalPages, pageSize)
	}

	return response
}

func (h *PaginationHelper) generateLinks(page, totalPages, pageSize int) *models.PaginationLinks {
	links := &models.PaginationLinks{
		First: h.buildURL(1, pageSize),
		Last:  h.buildURL(totalPages, pageSize),
	}

	if page < totalPages {
		next := h.buildURL(page+1, pageSize)
		links.Next = &next
	}

	if page > 1 {
		prev := h.buildURL(page-1, pageSize)
		links.Prev = &prev
	}

	return links
}

func (h *PaginationHelper) buildURL(page, pageSize int) string {
	return fmt.Sprintf("%s?page=%d&pageSize=%d", h.baseURL, page, pageSize)
}

func CalculateTotalPages(totalCount int64, pageSize int) int {
	if pageSize <= 0 {
		return 0
	}
	return int(math.Ceil(float64(totalCount) / float64(pageSize)))
}

func ValidatePaginationParams(page, pageSize int) (int, int) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}
	return page, pageSize
}

