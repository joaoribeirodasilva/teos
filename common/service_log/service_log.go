package service_log

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/joaoribeirodasilva/teos/common/database"
	"github.com/joaoribeirodasilva/teos/common/service_errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type LogEntry struct {
	ID          primitive.ObjectID `json:"_id"`
	Timestamp   time.Time          `json:"timestamp"`
	Type        string             `json:"type"`
	Application string             `json:"application"`
	Location    string             `json:"location"`
	Fields      string             `json:"fields,omitempty"`
	Message     string             `json:"message"`
	Code        int                `json:"code,omitempty"`
	HttpStatus  int                `json:"httpStatus,omitempty"`
}

var (
	ApplicationName string
	LogDatabase     *database.Db
	IsStdout        bool
	IsDatabase      bool
)

func InitServiceLog(appName string, database *database.Db) {
	ApplicationName = appName
	LogDatabase = database
}

func Info(location string, message string, args ...any) {
	logType("INFO", 0, 0, location, "", message, args...)
}

func Warning(location string, message string, args ...any) {
	logType("WARNNING", 0, 0, location, "", message, args...)
}

func Error(code int, httpCode int, location string, fields string, message string, args ...any) *service_errors.Error {
	args = append(args, fields)
	msg := logType("ERROR", 0, 0, location, fields, message, args)
	return service_errors.New(code, httpCode, location, fields, msg)
}

func Debug(location string, message string, args ...any) {
	logType("DEBUG", 0, 0, location, "", message, args...)
}

func logType(logType string, code int, httpCode int, location string, fields string, message string, args ...any) string {

	msg := makeMesage(location, message, args...)
	if IsStdout {
		slog.Info(msg)
	}
	if IsDatabase && LogDatabase != nil {

		entry := &LogEntry{
			ID:         primitive.NewObjectID(),
			Timestamp:  time.Now().UTC(),
			Type:       logType,
			Location:   location,
			Fields:     fields,
			Message:    msg,
			Code:       code,
			HttpStatus: httpCode,
		}

		logDb(entry)
	}
	return msg
}

func logDb(entry *LogEntry) error {

	coll := LogDatabase.Db.Collection("log_logs")
	if _, err := coll.InsertOne(context.TODO(), entry); err != nil {
		return err
	}
	return nil
}

func makeMesage(location string, message string, args ...any) string {
	msg := fmt.Sprintf(message, args...)
	return fmt.Sprintf("[%s] - %s", location, msg)
}
