package requests

import (
	"errors"
	"fmt"
	"net/url"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joaoribeirodasilva/teos/common/logger"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type QueryString struct {
	page     int64                `query:"p"`
	pageSize int64                `query:"ps"`
	all      bool                 `query:"a"`
	sort     string               `query:"s"`
	dir      string               `query:"d"`
	filter   string               `query:"f"`
	c        *gin.Context         `query:"-"`
	ID       *uint                `query:"-"`
	Options  *options.FindOptions `query:"-"`
	Related  bool                 `query:"r"`
	Filter   *string              `query:"-"`
}

const (
	defaultPage     = 0
	defaultPageSize = 10
	maxPageSize     = 100
	defaultDir      = "DESC"
	defaultSort     = "id"
	defaultAll      = false

	TYPE_NULL    = 0
	TYPE_STRING  = 1
	TYPE_NUMERIC = 2
	TYPE_BOOL    = 3
	TYPE_DATE    = 4

	OPERATION_EQUAL              = "$eq"
	OPERATION_NOTEQUAL           = "$ne"
	OPERATION_LESSTHAN           = "$lt"
	OPERATION_GREATERTHAN        = "$gt"
	OPERATION_LESSTHANOREQUAL    = "$lte"
	OPERATION_GREATERTHANOREQUAL = "$gte"
	OPERATION_BETWEEN            = "$bt"
	OPERATION_STARTSWITH         = "$sw"
	OPERATION_ENDSWITH           = "$ew"
	OPERATION_CONTAINS           = "$ct"
)

var (
	operations = map[string]string{
		OPERATION_EQUAL:              "=",
		OPERATION_NOTEQUAL:           "!=",
		OPERATION_LESSTHAN:           "<",
		OPERATION_GREATERTHAN:        ">",
		OPERATION_LESSTHANOREQUAL:    "<=",
		OPERATION_GREATERTHANOREQUAL: ">=",
		OPERATION_BETWEEN:            "`%s` BETWEEN%s",
		OPERATION_STARTSWITH:         "LIKE('%%%s')",
		OPERATION_ENDSWITH:           "LIKE('%s%%')",
		OPERATION_CONTAINS:           "LIKE('%%%s%%')",
	}
	nonFilterKeys = []string{"p", "ps", "a", "s", "d", "r"}
)

func NewQueryString(c *gin.Context) *QueryString {
	q := new(QueryString)
	q.page = defaultPage
	q.pageSize = defaultPageSize
	q.all = defaultAll
	q.sort = defaultSort
	q.dir = defaultDir
	q.filter = ""
	q.ID = nil
	q.Filter = nil
	q.Options = &options.FindOptions{}
	q.Related = false
	q.c = c
	return q
}

func (q *QueryString) Bind() *logger.HttpError {

	var err error

	strPage := q.c.DefaultQuery("p", "0")
	strPageSize := q.c.DefaultQuery("ps", strconv.Itoa(defaultPageSize))
	strAll := q.c.DefaultQuery("a", "false")
	strSort := q.c.DefaultQuery("s", "id")
	strDir := q.c.DefaultQuery("d", "DESC")
	strRelated := q.c.DefaultQuery("r", "false")
	//strFilter := q.c.Query("f")
	fullQuery := q.c.Request.URL.Query()

	//fmt.Printf("Query: %+v\n", fullQuery)

	strId := q.c.Params.ByName("id")

	//fmt.Printf("Begin: %+v\n", q)

	var tempInt int
	tempInt, err = strconv.Atoi(strPage)
	if err != nil {
		err := errors.New("invalid page requested (must be an integer)")
		fields := []string{"p"}
		return logger.Error(logger.LogStatusBadRequest, &fields, "", err, nil)
	}

	q.page = int64(tempInt)

	tempInt, err = strconv.Atoi(strPageSize)
	if err != nil {
		err := errors.New("invalid page size requested (must be an integer)")
		fields := []string{"ps"}
		return logger.Error(logger.LogStatusBadRequest, &fields, "", err, nil)
	}

	q.pageSize = int64(tempInt)

	q.all, err = strconv.ParseBool(strAll)
	if err != nil {
		err := errors.New("invalid all records requested (must be a boolean)")
		fields := []string{"a"}
		return logger.Error(logger.LogStatusBadRequest, &fields, "", err, nil)
	}

	q.Related, err = strconv.ParseBool(strRelated)
	if err != nil {
		q.Related = false
	}

	//fmt.Printf("%+v\n", query)

	if q.pageSize <= 0 {
		q.pageSize = defaultPageSize
	} else if q.pageSize > maxPageSize {
		q.pageSize = maxPageSize
	}

	q.sort = defaultSort
	strSort = strings.TrimSpace(strSort)
	if strSort != "" {
		q.sort = strSort
	}

	q.dir = strings.ToUpper(strDir)
	if q.dir != "ASC" && q.dir != "DESC" {
		err := errors.New("invalid sort direction, it can be only ASC or DESC")
		fields := []string{"d"}
		return logger.Error(logger.LogStatusBadRequest, &fields, "", err, nil)
	}

	if err := q.getID(strId); err != nil {
		err := errors.New("invalid record id")
		return logger.Error(logger.LogStatusBadRequest, nil, "", err, nil)
	}

	q.getPagination()
	q.getSort()

	if httpErr := q.parseFilter(fullQuery); httpErr != nil {
		return httpErr
	}

	//fmt.Printf("%+v\n", query)

	return nil
}

func (q *QueryString) parseFilter(filters url.Values) *logger.HttpError {

	finalFilter := ""

	if len(filters) == 0 {
		return nil
	}

	// by default we add the filter for records that are not deleted
	hasDeletedBy := url.Values.Get(filters, "deletedBy")
	hasDeletedAt := url.Values.Get(filters, "deletedAt")
	if hasDeletedBy == "" && hasDeletedAt == "" {
		finalFilter = "deleted_by IS NULL AND deleted_at IS NULL"
	}

	for key, filter := range filters {
		if slices.Contains(nonFilterKeys, key) {
			continue
		}

		fmt.Printf("Key: %s, Val: %+v\n", key, filter)
		params := strings.Split(filter[0], ",")
		if len(params) != 2 {
			err := fmt.Errorf("invalid filter '%s' parameter count", key)
			return logger.Error(logger.LogStatusBadRequest, nil, "", err, nil)
		}

		field := key
		if field == "" {
			err := errors.New("invalid filter field")
			return logger.Error(logger.LogStatusBadRequest, nil, "", err, nil)
		}

		value := params[1]
		operation, ok := operations[params[0]]
		if !ok {
			err := fmt.Errorf("invalid filter '%s' operation", key)
			return logger.Error(logger.LogStatusBadRequest, nil, "", err, nil)
		}

		values := strings.Split(value, "|")
		if len(values) != 1 && params[0] != OPERATION_BETWEEN {
			err := fmt.Errorf("invalid filter '%s' value", key)
			return logger.Error(logger.LogStatusBadRequest, nil, "", err, nil)
		} else if len(values) != 2 && params[1] == OPERATION_BETWEEN {
			err := fmt.Errorf("invalid filter '%s' between value", key)
			return logger.Error(logger.LogStatusBadRequest, nil, "", err, nil)
		}

		dataType := TYPE_STRING

		var valStart interface{}
		var valEnd interface{}

		valStart = values[0]

		if q.isNull(values[0]) {
			valStart = nil
			dataType = TYPE_NULL
		} else if dataType == TYPE_STRING {
			isBool := q.isBool(values[0])
			if isBool != nil {
				dataType = TYPE_BOOL
				valStart = *isBool
			}
		} else if dataType == TYPE_NUMERIC {
			number := q.isNumeric(values[0])
			if number != nil {
				dataType = TYPE_NUMERIC
				valStart = *number
			}
		} else if dataType == TYPE_DATE {
			date := q.isISODate(values[0])
			if date != nil {
				dataType = TYPE_DATE
				valStart = *date
			}
		}

		if len(values) == 2 {
			if dataType != TYPE_NULL {
				err := fmt.Errorf("filter '%s', a null can't be used with 2 values in the same operation", key)
				return logger.Error(logger.LogStatusBadRequest, nil, "", err, nil)
			}
			if dataType != TYPE_BOOL {
				err := fmt.Errorf("filter '%s', a boolean can't be used with 2 values in the same operation", key)
				return logger.Error(logger.LogStatusBadRequest, nil, "", err, nil)
			}
			if dataType != TYPE_BOOL {
				err := fmt.Errorf("filter '%s', a string can't be used with 2 values in the same operation", key)
				return logger.Error(logger.LogStatusBadRequest, nil, "", err, nil)
			}

			if operation == OPERATION_BETWEEN {

				date := q.isISODate(values[1])
				if date != nil {
					valEnd = *date
				} else if dataType != TYPE_DATE {
					err := fmt.Errorf("filter '%s', invalid second date value", key)
					return logger.Error(logger.LogStatusBadRequest, nil, "", err, nil)
				}

				number := q.isNumeric(values[1])
				if number != nil {
					valEnd = *number
				} else if dataType != TYPE_NUMERIC {
					err := fmt.Errorf("filter '%s', invalid second numeric value", key)
					return logger.Error(logger.LogStatusBadRequest, nil, "", err, nil)
				}

			} else {
				err := fmt.Errorf("filter '%s', use of 2 values outside a between operation", key)
				return logger.Error(logger.LogStatusBadRequest, nil, "", err, nil)
			}
		}

		var filt string
		switch params[1] {
		case OPERATION_BETWEEN:
			if dataType != TYPE_DATE && dataType != TYPE_NUMERIC {
				err := fmt.Errorf("filter '%s', between operations can only be used with data types numeric and date", key)
				return logger.Error(logger.LogStatusBadRequest, nil, "", err, nil)
			}
			if dataType != TYPE_DATE {
				startDate := valStart.(time.Time)
				endDate := valEnd.(time.Time)
				filt = fmt.Sprintf("('%s','%s') ", startDate.Format("2006-01-02 15:04:05"), endDate.Format("2006-01-02 15:04:05"))
			} else {
				filt = fmt.Sprintf("('%d','%d') ", valStart, valEnd)
			}
			filt = fmt.Sprintf(operation, field, filt)

		case OPERATION_STARTSWITH:
			if dataType != TYPE_STRING {
				err := fmt.Errorf("filter '%s', start with operations can only be used with data type string", key)
				return logger.Error(logger.LogStatusBadRequest, nil, "", err, nil)
			}
			filt = fmt.Sprintf(operation, valStart)
		case OPERATION_ENDSWITH:
			if dataType != TYPE_STRING {
				err := fmt.Errorf("filter '%s', end with operations can only be used with data type string", key)
				return logger.Error(logger.LogStatusBadRequest, nil, "", err, nil)
			}
			filt = fmt.Sprintf(operation, valStart)
		case OPERATION_CONTAINS:
			if dataType != TYPE_STRING {
				err := fmt.Errorf("filter '%s', contains operations can only be used with data type string", key)
				return logger.Error(logger.LogStatusBadRequest, nil, "", err, nil)
			}
			filt = fmt.Sprintf(operation, valStart)
		default:
			if dataType == TYPE_STRING {
				filt = fmt.Sprintf(" `%s`%s'%s' ", field, operation, valStart)
			} else if dataType == TYPE_NUMERIC {
				filt = fmt.Sprintf(" `%s`%s%d ", field, operation, valStart)
			} else if dataType == TYPE_BOOL {
				filt = fmt.Sprintf(" `%s`%s%t ", field, operation, valStart)
			} else if dataType == TYPE_DATE {
				startDate := valStart.(time.Time)
				filt = fmt.Sprintf(" `%s`%s%s ", field, operation, startDate.Format("2006-01-02 15:04:05"))
			} else {
				filt = fmt.Sprintf(" `%s`%s%s ", field, " IS ", valStart)
			}
		}

		finalFilter = fmt.Sprintf(finalFilter, filt)
	}

	q.Filter = &finalFilter

	fmt.Printf("Filter: %+v\n", q.Filter)

	return nil
}

func (q *QueryString) isNull(value string) bool {
	return strings.ToLower(value) == "null"
}

func (q *QueryString) isNumeric(value string) *float64 {

	val, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return nil
	}
	return &val
}

func (q *QueryString) isBool(value string) *bool {

	val, err := strconv.ParseBool(value)
	if err != nil {
		return nil
	}
	return &val
}

func (q *QueryString) isISODate(value string) *time.Time {
	t, err := time.Parse(time.RFC3339, value)
	if err != nil {
		return nil
	}
	return &t
}

func (q *QueryString) getID(strId string) error {

	id, err := strconv.Atoi(strId)
	if err != nil {
		return err
	}
	final := uint(id)
	q.ID = &final
	return nil
}

func (q *QueryString) getPagination() {

	if q.all {
		q.Options.SetSkip(q.pageSize * q.page)
	} else {
		q.Options.SetLimit(int64(q.pageSize))
		q.Options.SetSkip(q.pageSize * q.page)
	}
}

func (q *QueryString) getSort() {
	q.Options.SetSort(bson.D{{Key: q.sort, Value: q.dir}})
}
