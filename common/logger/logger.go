package logger

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/joaoribeirodasilva/teos/common/database"
	"github.com/joaoribeirodasilva/teos/common/utils/dump"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type LogStatus uint8
type LogTypes uint8
type LogLevel uint8

const (
	LogStatusNone LogStatus = iota
	LogStatusNotFound
	LogStatusBadRequest
	LogStatusConflict
	LogStatusInternalServerError
	LogStatusUnauthorized
	LogStatusForbidden
)

const (
	LogTypesError LogTypes = iota
	LogTypesWarn
	LogTypesInfo
	LogTypesDebug
)

const (
	LogLevelNone LogLevel = iota
	LogLevelError
	LogLevelWarning
	LogLevelInfo
	LogLevelDebug
)

var (
	typeNames = map[LogTypes]string{
		LogTypesError: "ERROR",
		LogTypesWarn:  "WARN",
		LogTypesInfo:  "INFO",
		LogTypesDebug: "DEBUG",
	}

	httpStatus = map[LogStatus]int{
		LogStatusNotFound:            404,
		LogStatusBadRequest:          400,
		LogStatusConflict:            409,
		LogStatusUnauthorized:        401,
		LogStatusForbidden:           403,
		LogStatusInternalServerError: 500,
	}

	Logger Log
)

type Log struct {
	Db             *database.Db
	Application    string
	CollectionName string
	LogLevel       LogLevel
}

type LogDocument struct {
	ID          primitive.ObjectID `json:"_id" bson:"_id"`
	Application string             `json:"application" bson:"application"`
	Timestamp   time.Time          `json:"timestamp" bson:"timestamp"`
	Message     *LogMessage        `json:"message" bson:"message"`
}

type LogMessage struct {
	Type    string    `json:"type" bson:"type"`
	Fields  *[]string `json:"fields,omitempty" bson:"fields,omitempty"`
	Message string    `json:"message" bson:"message"`
	Data    *string   `json:"data,omitempty" bson:"data,omitempty"`
}

type HttpError struct {
	Status LogStatus  `json:"-"`
	Err    LogMessage `json:"error"`
}

func (e *HttpError) Error() string {

	return e.Err.Message
}

func Info(message string, args ...any) {

	if Logger.LogLevel < LogLevelInfo {
		return
	}

	txt := fmt.Sprintf(message, args...)
	msg := createMessage(LogTypesInfo, txt, nil, nil)

	slog.Info(msg.Message)

}

func Warn(message string, args ...any) {

	if Logger.LogLevel < LogLevelWarning {
		return
	}

	txt := fmt.Sprintf(message, args...)
	msg := createMessage(LogTypesWarn, txt, nil, nil)

	slog.Warn(msg.Message)

}

func Error(status LogStatus, fields *[]string, message string, err error, data any) *HttpError {

	httpErr := &HttpError{
		Status: LogStatus(httpStatus[status]),
		Err:    LogMessage{},
	}

	txt := ""
	if err != nil {

		txt = fmt.Sprintf("%s. ERR: %s", message, err)

	} else {

		txt = fmt.Sprintf("%s.", message)
	}

	msg := createMessage(LogTypesError, txt, fields, data)

	httpErr.Err = *msg

	if Logger.LogLevel < LogLevelError {
		return httpErr
	}

	slog.Error(msg.Message)
	if msg.Data != nil {
		slog.Error(fmt.Sprintf("Data: %s", *msg.Data))
	}

	return httpErr
}

func Debug(message string, data any, args ...any) {

	if Logger.LogLevel < LogLevelDebug {
		return
	}

	txt := fmt.Sprintf(message, args...)
	msg := createMessage(LogTypesDebug, txt, nil, data)

	slog.Debug(msg.Message)
	if msg.Data != nil {
		slog.Debug(*msg.Data)
	}

}

func DebugDump(obj interface{}) {

	slog.Debug(dump.ToJSON(obj))

}

func SetApplication(application string) {
	Logger.Application = application
}

func GetApplication() string {
	return Logger.Application
}

func SetDatabase(db *database.Db) {
	Logger.Db = db
}

func SetCollectionName(collectionName string) {
	Logger.CollectionName = collectionName
}

func SetLogLevel(level LogLevel) {
	Logger.LogLevel = level
}

func createMessage(typ LogTypes, mesg string, fields *[]string, data any) *LogMessage {

	var dat *string = nil
	if data != nil {
		d, _ := dump.ToJSON(data)
		dat = &d
	}

	msg := LogMessage{
		Type:    typeNames[typ],
		Fields:  fields,
		Message: mesg,
		Data:    dat,
	}

	persist(msg)

	return &msg
}

func persist(msg LogMessage) {

	if Logger.Db == nil {
		return
	}

	if Logger.CollectionName == "" {
		slog.Warn("no collection name set for Logger, log persistance disabled")
		return
	}

	coll := Logger.Db.GetDatabase().Collection(Logger.CollectionName)

	doc := &LogDocument{
		ID:          primitive.NewObjectID(),
		Application: Logger.Application,
		Timestamp:   time.Now().UTC(),
		Message:     &msg,
	}

	if _, err := coll.InsertOne(Logger.Db.GetContext(), doc); err != nil {
		slog.Error(fmt.Sprintf("failed to persist log message, ERR: %s", err.Error()))
	}
}
