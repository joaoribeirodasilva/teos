package requests

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joaoribeirodasilva/teos/common/service_errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type QueryString struct {
	page     int64                `query:"p"`
	pageSize int64                `query:"ps"`
	all      bool                 `query:"a"`
	sort     string               `query:"s"`
	dir      int                  `query:"d"`
	filter   string               `query:"f"`
	c        *gin.Context         `query:"-"`
	ID       *primitive.ObjectID  `query:"-"`
	Options  *options.FindOptions `query:"-"`
	Filter   *primitive.D         `query:"-"`
}

const (
	defaultPage     = 0
	defaultPageSize = 10
	maxPageSize     = 100
	defaultDir      = -1
	defaultSort     = "_id"
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

var operations = map[string]string{
	OPERATION_EQUAL:              "$eq",
	OPERATION_NOTEQUAL:           "$ne",
	OPERATION_LESSTHAN:           "$lt",
	OPERATION_GREATERTHAN:        "$gt",
	OPERATION_LESSTHANOREQUAL:    "$lte",
	OPERATION_GREATERTHANOREQUAL: "$gte",
	OPERATION_BETWEEN:            "",
	OPERATION_STARTSWITH:         "",
	OPERATION_ENDSWITH:           "",
	OPERATION_CONTAINS:           "",
}

func NewQueryString(c *gin.Context) *QueryString {
	q := new(QueryString)
	q.page = defaultPage
	q.pageSize = defaultPageSize
	q.all = defaultAll
	q.sort = defaultSort
	q.dir = defaultDir
	q.filter = ""
	q.ID = nil
	q.Filter = &primitive.D{{}}
	q.Options = &options.FindOptions{}
	q.c = c
	return q
}

func (q *QueryString) Bind() *service_errors.Error {

	var err error

	strPage := q.c.DefaultQuery("p", "0")
	strPageSize := q.c.DefaultQuery("ps", strconv.Itoa(defaultPageSize))
	strAll := q.c.DefaultQuery("a", "false")
	strSort := q.c.DefaultQuery("s", "id")
	strDir := q.c.DefaultQuery("d", "DESC")
	strFilter := q.c.Query("f")
	strId := q.c.Params.ByName("id")

	//fmt.Printf("Begin: %+v\n", q)

	var tempInt int
	tempInt, err = strconv.Atoi(strPage)
	if err != nil {
		return service_errors.New(0, http.StatusBadRequest, "REQUESTS", "Bind", "invalid page requested (must be an integer)", "").LogError()
	}

	q.page = int64(tempInt)

	tempInt, err = strconv.Atoi(strPageSize)
	if err != nil {
		return service_errors.New(0, http.StatusBadRequest, "REQUESTS", "Bind", "invalid page size requested (must be an integer)", "").LogError()
	}
	q.pageSize = int64(tempInt)

	q.all, err = strconv.ParseBool(strAll)
	if err != nil {
		return service_errors.New(0, http.StatusBadRequest, "REQUESTS", "Bind", "invalid all records requested (must be a boolean)", "").LogError()
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

	q.dir = 1
	if q.sort == "_id" {
		q.dir = -1
	}

	strDir = strings.ToUpper(strDir)
	if strDir == "DESC" {
		q.dir = -1
	} else if strDir == "ASC" {
		q.dir = 1
	}

	if err := q.getID(strId); err != nil {
		service_errors.New(0, http.StatusBadRequest, "REQUESTS", "Bind", "invalid record id. ERR: %s", err.Error()).LogError()
	}

	q.getPagination()
	q.getSort()

	if appErr := q.parseFilter(strFilter); appErr != nil {
		return appErr
	}

	//fmt.Printf("%+v\n", query)

	return nil
}

func (q *QueryString) parseFilter(filter string) *service_errors.Error {

	filterArray := []primitive.D{}

	if filter == "" {
		return nil
	}

	filters := strings.Split(filter, ";")

	for _, filter := range filters {

		params := strings.Split(filter, ",")
		if len(params) != 3 {
			return service_errors.New(0, http.StatusBadRequest, "REQUESTS", "parseFilter", "invalid filter parameter count", "").LogError()
		}

		field := params[0]
		if field == "" {
			return service_errors.New(0, http.StatusBadRequest, "REQUESTS", "parseFilter", "invalid filter field", "").LogError()
		}

		value := params[2]
		operation, ok := operations[params[1]]
		if !ok {
			return service_errors.New(0, http.StatusBadRequest, "REQUESTS", "parseFilter", "invalid filter operation", "").LogError()
		}

		values := strings.Split(value, "|")
		if len(values) != 1 && params[1] != OPERATION_BETWEEN {
			return service_errors.New(0, http.StatusBadRequest, "REQUESTS", "parseFilter", "invalid filter value", "").LogError()
		} else if len(values) != 2 && params[1] == OPERATION_BETWEEN {
			return service_errors.New(0, http.StatusBadRequest, "REQUESTS", "parseFilter", "invalid filter between value", "").LogError()
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
		} else if dataType == TYPE_STRING {
			number := q.isNumeric(values[0])
			if number != nil {
				dataType = TYPE_NUMERIC
				valStart = *number
			}
		} else if dataType == TYPE_STRING {
			date := q.isISODate(values[0])
			if date != nil {
				dataType = TYPE_DATE
				valStart = *date
			}
		}

		if len(values) == 2 {
			if dataType != TYPE_NULL {
				return service_errors.New(0, http.StatusBadRequest, "REQUESTS", "parseFilter", "a null can't be used with 2 values in the same operation", "").LogError()
			}
			if dataType != TYPE_BOOL {
				return service_errors.New(0, http.StatusBadRequest, "REQUESTS", "parseFilter", "a boolean can't be used with 2 values in the same operation", "").LogError()
			}
			if dataType != TYPE_BOOL {
				return service_errors.New(0, http.StatusBadRequest, "REQUESTS", "parseFilter", "a string can't be used with 2 values in the same operation", "").LogError()
			}

			if operation == OPERATION_BETWEEN {

				date := q.isISODate(values[1])
				if date != nil {
					valEnd = *date
				} else if dataType != TYPE_DATE {
					return service_errors.New(0, http.StatusBadRequest, "REQUESTS", "parseFilter", "invalid second date value", "").LogError()
				}

				number := q.isNumeric(values[1])
				if number != nil {
					valEnd = *number
				} else if dataType != TYPE_NUMERIC {
					return service_errors.New(0, http.StatusBadRequest, "REQUESTS", "parseFilter", "invalid second numeric value", "").LogError()
				}

			} else {
				return service_errors.New(0, http.StatusBadRequest, "REQUESTS", "parseFilter", "use of 2 values outside a between operation", "").LogError()
			}
		}

		var filt primitive.D
		switch params[1] {
		case OPERATION_BETWEEN:
			if dataType != TYPE_DATE && dataType != TYPE_NUMERIC {
				return service_errors.New(0, http.StatusBadRequest, "REQUESTS", "parseFilter", "between operations can only be used with data types numeric and date", "").LogError()
			}
			filt = primitive.D{
				{Key: "$and", Value: primitive.A{
					primitive.D{{Key: field, Value: primitive.D{{Key: "$gte", Value: valStart}}}},
					primitive.D{{Key: field, Value: primitive.D{{Key: "$lte", Value: valEnd}}}},
				}},
			}
		case OPERATION_STARTSWITH:
			if dataType != TYPE_STRING {
				return service_errors.New(0, http.StatusBadRequest, "REQUESTS", "parseFilter", "start with operations can only be used with data type string", "").LogError()
			}
			filt = primitive.D{{Key: field, Value: primitive.D{{Key: "$regex", Value: primitive.Regex{Pattern: "^" + values[0], Options: "i"}}}}}
		case OPERATION_ENDSWITH:
			if dataType != TYPE_STRING {
				return service_errors.New(0, http.StatusBadRequest, "REQUESTS", "parseFilter", "end with operations can only be used with data type string", "").LogError()
			}
			filt = primitive.D{{Key: field, Value: primitive.D{{Key: "$regex", Value: primitive.Regex{Pattern: values[0] + "$", Options: "i"}}}}}
		case OPERATION_CONTAINS:
			if dataType != TYPE_STRING {
				return service_errors.New(0, http.StatusBadRequest, "REQUESTS", "parseFilter", "contains operations can only be used with data type string", "").LogError()
			}
			filt = primitive.D{{Key: field, Value: primitive.D{{Key: "$regex", Value: primitive.Regex{Pattern: values[0], Options: "i"}}}}}
		default:
			filt = primitive.D{{Key: field, Value: primitive.D{{Key: operation, Value: valStart}}}}
		}

		filterArray = append(filterArray, filt)
	}

	if len(filterArray) > 1 {
		q.Filter = &primitive.D{{Key: "$and", Value: filterArray}}
	} else if len(filterArray) == 1 {
		q.Filter = &filterArray[0]
	}

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

	id, err := primitive.ObjectIDFromHex(strId)
	if err != nil {
		return err
	}
	q.ID = &id
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
