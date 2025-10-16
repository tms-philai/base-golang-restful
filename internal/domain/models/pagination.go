package models

type PaginationRequest struct {
	Page     int    `form:"page" json:"page" binding:"min=1"`
	PageSize int    `form:"pageSize" json:"pageSize" binding:"min=1,max=100"`
	SortBy   string `form:"sortBy" json:"sortBy"`
	Order    string `form:"order" json:"order" binding:"omitempty,oneof=asc desc ASC DESC"`
}

type PaginationResponse struct {
	Data       interface{}      `json:"data"`
	Pagination PaginationMeta   `json:"pagination"`
	Links      *PaginationLinks `json:"links,omitempty"`
}

type PaginationMeta struct {
	Page       int   `json:"page"`
	PageSize   int   `json:"pageSize"`
	TotalCount int64 `json:"totalCount"`
	TotalPages int   `json:"totalPages"`
	HasNext    bool  `json:"hasNext"`
	HasPrev    bool  `json:"hasPrev"`
}

type PaginationLinks struct {
	First string  `json:"first"`
	Last  string  `json:"last"`
	Next  *string `json:"next"`
	Prev  *string `json:"prev"`
}

type CursorPaginationRequest struct {
	Cursor   string `form:"cursor" json:"cursor"`
	PageSize int    `form:"pageSize" json:"pageSize" binding:"min=1,max=100"`
	SortBy   string `form:"sortBy" json:"sortBy"`
	Order    string `form:"order" json:"order" binding:"omitempty,oneof=asc desc ASC DESC"`
}

type CursorPaginationResponse struct {
	Data       interface{}          `json:"data"`
	Pagination CursorPaginationMeta `json:"pagination"`
}

type CursorPaginationMeta struct {
	NextCursor *string `json:"nextCursor"`
	PrevCursor *string `json:"prevCursor"`
	HasNext    bool    `json:"hasNext"`
	HasPrev    bool    `json:"hasPrev"`
	PageSize   int     `json:"pageSize"`
}

func NewPaginationRequest() *PaginationRequest {
	return &PaginationRequest{
		Page:     1,
		PageSize: 10,
		Order:    "asc",
	}
}

func (p *PaginationRequest) Validate() {
	if p.Page < 1 {
		p.Page = 1
	}
	if p.PageSize < 1 {
		p.PageSize = 10
	}
	if p.PageSize > 100 {
		p.PageSize = 100
	}
	if p.Order == "" {
		p.Order = "asc"
	}
}

func (p *PaginationRequest) GetOffset() int {
	return (p.Page - 1) * p.PageSize
}

func (p *PaginationRequest) GetLimit() int {
	return p.PageSize
}

func NewCursorPaginationRequest() *CursorPaginationRequest {
	return &CursorPaginationRequest{
		PageSize: 10,
		Order:    "asc",
	}
}

func (c *CursorPaginationRequest) Validate() {
	if c.PageSize < 1 {
		c.PageSize = 10
	}
	if c.PageSize > 100 {
		c.PageSize = 100
	}
	if c.Order == "" {
		c.Order = "asc"
	}
}
