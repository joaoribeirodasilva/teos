package payload

import (
	"strings"

	"github.com/gin-gonic/gin"
)

type SortOrder string

const (
	SORT_ASC  SortOrder = "asc"
	SORT_DESC SortOrder = "desc"

	DEFAULT_SORT_FIELD = "updated_at"
	DEFAULT_SORT_DIR   = SORT_DESC
)

type Sort struct {
	ctx   *gin.Context
	Field string
	Order string
}

func NewSort(ctx *gin.Context) *Sort {
	// DEFAULTS
	return &Sort{
		ctx:   ctx,
		Field: DEFAULT_SORT_FIELD,
		Order: string(DEFAULT_SORT_DIR),
	}
}

func (s *Sort) Parse() error {

	s.Field = s.ctx.DefaultQuery("s", DEFAULT_SORT_FIELD)
	s.Order = s.ctx.DefaultQuery("d", string(DEFAULT_SORT_DIR))
	s.Order = strings.ToLower(s.Order)
	if s.Order != string(SORT_DESC) && s.Order != string(SORT_ASC) {
		s.Order = string(DEFAULT_SORT_DIR)
	}

	return nil
}
