package payload

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

const (
	DEFAULT_PAGE      = 0
	DEFAULT_PAGE_SIZE = 20
	DEFAULT_PAGE_ALL  = false
)

type Pagination struct {
	ctx      *gin.Context
	Page     uint
	PageSize uint
	All      bool
}

func NewPagination(ctx *gin.Context) *Pagination {

	return &Pagination{
		ctx:      ctx,
		Page:     DEFAULT_PAGE,
		PageSize: DEFAULT_PAGE_SIZE,
		All:      DEFAULT_PAGE_ALL,
	}
}

func (p *Pagination) Parse() error {

	page, err := strconv.Atoi(p.ctx.DefaultQuery("p", "0"))
	p.Page = uint(page)
	if err != nil {
		return err
	}

	pageSize, err := strconv.Atoi(p.ctx.DefaultQuery("ps", "0"))
	p.PageSize = uint(pageSize)
	if err != nil {
		return err
	}

	all := p.ctx.DefaultQuery("a", "false")
	p.All, err = strconv.ParseBool(all)
	if err != nil {
		return err
	}

	return nil
}
