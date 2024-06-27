package payload

import (
	"fmt"
	"net/url"
	"slices"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	FILTER_EQ  FilterOperation = "eq"
	FILTER_NE  FilterOperation = "ne"
	FILTER_GT  FilterOperation = "gt"
	FILTER_GTE FilterOperation = "gte"
	FILTER_LT  FilterOperation = "lt"
	FILTER_LTE FilterOperation = "lte"
	FILTER_BT  FilterOperation = "bt"
	FILTER_SW  FilterOperation = "sw"
	FILTER_EW  FilterOperation = "ew"
	FILTER_CT  FilterOperation = "ct"
)

var (
	formatOperator = map[FilterOperation]string{
		FILTER_EQ:  "`%s` = '%s'",
		FILTER_NE:  "`%s` != '%s'",
		FILTER_GT:  "`%s` > '%s'",
		FILTER_GTE: "`%s` >= '%s'",
		FILTER_LT:  "`%s` < '%s'",
		FILTER_LTE: "`%s` <= '%s'",
		FILTER_BT:  "`%s`'%s' BETWEEN '%s' AND '%s'",
		FILTER_SW:  "`%s`LIKE('%%%s')",
		FILTER_EW:  "`%s`LIKE('%s%%')",
		FILTER_CT:  "`%s`LIKE('%%%s%%')",
	}

	validOperations = []FilterOperation{
		FILTER_EQ,
		FILTER_NE,
		FILTER_GT,
		FILTER_GTE,
		FILTER_LT,
		FILTER_LTE,
		FILTER_BT,
		FILTER_SW,
		FILTER_EW,
		FILTER_CT,
	}
	nonFilterKeys = []string{"p", "ps", "a", "s", "d", "r"}
)

type FilterOperation string

type Filter struct {
	ctx    *gin.Context
	Filter string
}

func NewFilter(ctx *gin.Context) *Filter {

	return &Filter{
		ctx: ctx,
	}
}

func (f *Filter) Parse() error {

	query := f.ctx.Request.URL.Query()

	if !query.Has("deleted_by") && !query.Has("deleted_at") {
		query.Add("deleted_at", "eq,NULL")
		query.Add("deleted_by", "eq,NULL")
	}

	if err := f.parseValues(query); err != nil {
		return err
	}

	return nil
}

func (f *Filter) parseValues(values url.Values) error {

	if len(values) == 0 {
		return nil
	}

	for key, val := range values {

		if slices.Contains(nonFilterKeys, key) {
			continue
		}

		operationAndValue := strings.Split(val[0], ",")
		if len(operationAndValue) != 2 {
			return fmt.Errorf("invalid operation or value for filter %s", key)
		}

		operation := operationAndValue[0]
		filterValues := operationAndValue[1]

		if !f.isValidOperation(operation) {
			return fmt.Errorf("invalid operation %s for filter %s", operation, key)
		}

		if filterValues == "" {
			return fmt.Errorf("invalid empty value for filter %s", key)
		}

		startValue := ""
		endValue := ""
		allValues := strings.Split(filterValues, "|")
		numValues := 1
		if len(allValues) == 2 {
			numValues = 2
			endValue = allValues[1]
		}
		startValue = allValues[0]

		if (FilterOperation(operation) != FILTER_BT && numValues == 2) || FilterOperation(operation) == FILTER_BT && numValues != 2 {
			return fmt.Errorf("invalid value for operation %s in filter %s", operation, key)
		}

		filter := ""
		format := formatOperator[FilterOperation(operation)]
		if len(allValues) == 2 {
			filter = fmt.Sprintf(format, key, startValue, endValue)
		} else {
			filter = fmt.Sprintf(format, key, startValue)
		}

		if strings.ToUpper(startValue) == "NULL" {
			filter = fmt.Sprintf("`%s` IS NULL", key)
		}

		if f.Filter != "" {
			f.Filter += " AND "
		}

		f.Filter += filter
	}

	return nil
}

func (f *Filter) isValidOperation(operation string) bool {
	return slices.Contains(validOperations, FilterOperation(operation))
}
