package payload

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

type HttpQuery struct {
	ctx        *gin.Context
	Pagination *Pagination
	Sort       *Sort
	Filter     *Filter
	Related    bool
}

func NewHttpQuery(ctx *gin.Context) *HttpQuery {
	return &HttpQuery{
		ctx:        ctx,
		Pagination: NewPagination(ctx),
		Sort:       NewSort(ctx),
		Filter:     NewFilter(ctx),
	}
}

func (h *HttpQuery) Parse() error {

	var err error

	related := h.ctx.DefaultQuery("r", "false")
	h.Related, err = strconv.ParseBool(related)
	if err != nil {
		return err
	}

	if err := h.Pagination.Parse(); err != nil {
		return err
	}
	if err := h.Sort.Parse(); err != nil {
		return err
	}
	if err := h.Filter.Parse(); err != nil {
		return err
	}

	return nil
}
