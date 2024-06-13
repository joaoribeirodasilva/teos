package service_log

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/joaoribeirodasilva/teos/common/redisdb"
	"github.com/joaoribeirodasilva/teos/common/service_errors"
)

type LogEntry struct {
	EntryID     string    `json:"entryId"`
	Timestamp   time.Time `json:"timestamp"`
	Type        string    `json:"type"`
	Application string    `json:"application"`
	Location    string    `json:"location"`
	Fields      string    `json:"fields,omitempty"`
	Message     string    `json:"message"`
	Code        int       `json:"code,omitempty"`
	HttpStatus  int       `json:"httpStatus,omitempty"`
}

var (
	ApplicationName string
	LogDatabase     *redisdb.RedisDB
	IsStdout        bool
	IsDatabase      bool
)

func InitServiceLog(addr string, db int, username string, password string) error {
	LogDatabase := redisdb.New("Logs Database", addr, db, username, password)
	if err := LogDatabase.Connect(); err != nil {
		return err
	}
	return nil
}

func Info(location string, message string, args ...string) {
	logType("INFO", 0, 0, location, "", message, args)
}

func Warning(location string, message string, args ...string) {
	logType("WARNNING", 0, 0, location, "", message, args)
}

func Error(code int, httpCode int, location string, fields string, message string, args ...any) *service_errors.Error {
	msg := logType("ERROR", 0, 0, location, fields, message, args)
	return service_errors.New(code, httpCode, location, fields, msg)
}

func Debug(location string, message string, args ...string) {
	logType("DEBUG", 0, 0, location, "", message, args)
}

func logType(logType string, code int, httpCode int, location string, fields string, message string, args ...any) string {

	msg := makeMesage(location, message, args)
	if IsStdout {
		slog.Info(msg)
	}
	if IsDatabase && LogDatabase != nil {

		entry := &LogEntry{
			EntryID:    uuid.NewString(),
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

	key := fmt.Sprintf("log_%s", entry.EntryID)
	if err := LogDatabase.Client.Set(context.Background(), key, entry, 0); err != nil {
		return err.Err()
	}
	return nil
}

func makeMesage(message string, args ...any) string {
	return fmt.Sprintf("[%s] - "+message, args)
}
