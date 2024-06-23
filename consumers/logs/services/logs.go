package services

import (
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

	/* 	models, err := s.loadData()
	   	if err != nil {
	   		return nil, err
	   	} */

	return nil
}
