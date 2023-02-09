package common

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Pagination struct {
	Metadata struct {
		Limit      int    `json:"limit,omitempty;query:limit"`
		Page       int    `json:"page,omitempty;query:page"`
		Sort       string `json:"sort,omitempty;query:sort"`
		TotalRows  int64  `json:"total_rows"`
		TotalPages int    `json:"total_pages"`
	} `json:"metadata"`
	Data interface{} `json:"data"`
}

func (p *Pagination) GetOffset() int {
	return (p.GetPage() - 1) * p.GetLimit()
}

func (p *Pagination) GetLimit() int {
	if p.Metadata.Limit <= 0 {
		p.Metadata.Limit = 10
	}
	return p.Metadata.Limit
}

func (p *Pagination) GetPage() int {
	if p.Metadata.Page <= 0 {
		p.Metadata.Page = 1
	}
	return p.Metadata.Page
}

func (p *Pagination) GetSort() string {
	if p.Metadata.Sort == "" {
		p.Metadata.Sort = "id desc"
	}
	return p.Metadata.Sort
}

func GetPaginationParameter(c *gin.Context, p *Pagination) {
	page := c.Query("page")
	if page == "" {
		page = "0"
	}
	p.Metadata.Page, _ = strconv.Atoi(page)

	limit := c.Query("limit")
	if limit == "" {
		limit = "0"
	}
	p.Metadata.Limit, _ = strconv.Atoi(limit)

	order := strings.ToLower(c.Query("order"))
	if order == "" {
		order = "desc"
	}

	sort := strings.ToLower(c.Query("sort"))
	if sort == "" {
		sort = "id"
	}
	p.Metadata.Sort = sort + " " + order
	// fmt.Println(p.Page, p.Limit, p.Sort)
}

func Paginate(value interface{}, pagination *Pagination, db *gorm.DB) func(db *gorm.DB) *gorm.DB {
	var totalRows int64
	db.Model(value).Count(&totalRows)

	pagination.Metadata.TotalRows = totalRows
	totalPages := int(math.Ceil(float64(totalRows) / float64(pagination.GetLimit())))
	fmt.Println(totalPages)
	pagination.Metadata.TotalPages = totalPages
	fmt.Println(pagination)

	return func(db *gorm.DB) *gorm.DB {
		return db.Offset(pagination.GetOffset()).Limit(pagination.GetLimit()).Order(pagination.GetSort())
	}
}

func NewPagination() Pagination {
	var pagination Pagination
	return pagination
}
