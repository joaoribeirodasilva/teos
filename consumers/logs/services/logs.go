package services

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/joaoribeirodasilva/teos/common/logger"
	"github.com/joaoribeirodasilva/teos/common/models"
	"github.com/joaoribeirodasilva/teos/common/payload"
	"github.com/joaoribeirodasilva/teos/common/redisdb"
	"gorm.io/gorm"
)

type LogStats struct {
	Millis   int64
	Count    int64
	Interval int64
}

type LogsService struct {
	payload *payload.Payload
	db      *gorm.DB
	request *payload.HttpRequest
	history *redisdb.RedisDB
}

func NewLogsService(payload *payload.Payload) *LogsService {
	return &LogsService{
		payload: payload,
		db:      payload.Services.Db.GetDatabase(),
		request: payload.Http.Request,
		history: payload.Services.HistoryDb,
	}
}

func (s *LogsService) List() (*models.LogMessages, *logger.HttpError) {

	/* 	models, err := s.loadData()
	   	if err != nil {
	   		return nil, err
	   	} */

	return nil, nil
}

func (s *LogsService) Get(*models.LogMessage) *logger.HttpError {

	/* 	models, err := s.loadData()
	   	if err != nil {
	   		return nil, err
	   	} */

	return nil
}

func (s *LogsService) Stats() *logger.HttpError {

	/* 	models, err := s.loadData()
	   	if err != nil {
	   		return nil, err
	   	} */

	return nil
}

func (s *LogsService) Save() *logger.HttpError {

	logRecords, err := s.payload.Services.LogsDb.MGet("*")
	if err != nil {
		slog.Warn(fmt.Sprintf("failed to read logs from memory database, ERR: %s", err.Error()))
		return nil
	}

	if len(*logRecords) == 0 {
		return nil
	}

	for k, v := range *logRecords {

		if v == nil {
			continue
		}

		jsonModel := logger.LogMessage{}
		if err := json.Unmarshal([]byte(*v), &jsonModel); err != nil {
			slog.Warn(fmt.Sprintf("failed to decode json from log entry, ERR: %s", err.Error()))
			continue
		}

		now := time.Now().UTC()
		logModel := models.LogMessage{
			Time:      jsonModel.Time,
			App:       jsonModel.App,
			Type:      jsonModel.Type,
			Data:      v,
			CreatedBy: 1,
			CreatedAt: now,
		}

		if err := s.payload.Services.Db.GetDatabase().Model(&logModel).Create(&logModel).Error; err != nil {
			slog.Warn(fmt.Sprintf("failed to save log entry into the database, ERR: %s", err.Error()))
			continue
		}

		if err := s.payload.Services.LogsDb.Del(k); err != nil {
			slog.Warn(fmt.Sprintf("failed to delete log entry from memory database, ERR: %s", err.Error()))
		}

	}

	return nil
}
