package utils

import (
	"base-golang-restful-app/models"
	"fmt"
	"strings"

	"gorm.io/gorm"
)

type QueryBuilder struct {
	db *gorm.DB
}

func NewQueryBuilder(db *gorm.DB) *QueryBuilder {
	return &QueryBuilder{db: db}
}

func (qb *QueryBuilder) ApplyPagination(query *gorm.DB, req *models.PaginationRequest) *gorm.DB {
	req.Validate()

	if req.SortBy != "" {
		order := fmt.Sprintf("%s %s", req.SortBy, strings.ToUpper(req.Order))
		query = query.Order(order)
	}

	offset := req.GetOffset()
	limit := req.GetLimit()

	return query.Offset(offset).Limit(limit)
}

func (qb *QueryBuilder) ApplySearch(query *gorm.DB, searchFields []string, searchTerm string) *gorm.DB {
	if searchTerm == "" || len(searchFields) == 0 {
		return query
	}

	searchPattern := "%" + searchTerm + "%"
	conditions := make([]string, 0, len(searchFields))
	args := make([]interface{}, 0, len(searchFields))

	for _, field := range searchFields {
		conditions = append(conditions, fmt.Sprintf("%s ILIKE ?", field))
		args = append(args, searchPattern)
	}

	whereClause := strings.Join(conditions, " OR ")
	return query.Where(whereClause, args...)
}

func (qb *QueryBuilder) ApplyFilters(query *gorm.DB, filters map[string]interface{}) *gorm.DB {
	for key, value := range filters {
		if value == nil {
			query = query.Where(fmt.Sprintf("%s IS NULL", key))
		} else {
			query = query.Where(fmt.Sprintf("%s = ?", key), value)
		}
	}
	return query
}

func (qb *QueryBuilder) ApplyDateRange(query *gorm.DB, field string, start, end interface{}) *gorm.DB {
	if start != nil && end != nil {
		query = query.Where(fmt.Sprintf("%s BETWEEN ? AND ?", field), start, end)
	} else if start != nil {
		query = query.Where(fmt.Sprintf("%s >= ?", field), start)
	} else if end != nil {
		query = query.Where(fmt.Sprintf("%s <= ?", field), end)
	}
	return query
}

func (qb *QueryBuilder) ApplyMultiSort(query *gorm.DB, sorts []SortField) *gorm.DB {
	for _, sort := range sorts {
		order := fmt.Sprintf("%s %s", sort.Field, strings.ToUpper(sort.Order))
		query = query.Order(order)
	}
	return query
}

type SortField struct {
	Field string
	Order string
}

func ParseSortFields(sortBy, order string) []SortField {
	if sortBy == "" {
		return nil
	}

	fields := strings.Split(sortBy, ",")
	orders := strings.Split(order, ",")

	sorts := make([]SortField, 0, len(fields))
	for i, field := range fields {
		field = strings.TrimSpace(field)
		if field == "" {
			continue
		}

		sortOrder := "asc"
		if i < len(orders) && orders[i] != "" {
			sortOrder = strings.ToLower(strings.TrimSpace(orders[i]))
		}

		sorts = append(sorts, SortField{
			Field: field,
			Order: sortOrder,
		})
	}

	return sorts
}
