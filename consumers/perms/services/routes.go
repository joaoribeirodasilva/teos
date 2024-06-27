package services

import (
	"errors"

	"github.com/joaoribeirodasilva/teos/common/logger"
	"github.com/joaoribeirodasilva/teos/common/models"
	"github.com/joaoribeirodasilva/teos/common/payload"
	"github.com/joaoribeirodasilva/teos/common/redisdb"
	"github.com/joaoribeirodasilva/teos/common/utils/permission_key"

	"gorm.io/gorm"
)

type RoutesService struct {
	payload *payload.Payload
	db      *gorm.DB
	request *payload.HttpRequest
	history *redisdb.RedisDB
}

func NewRoutesService(payload *payload.Payload) *RoutesService {
	return &RoutesService{
		payload: payload,
		db:      payload.Services.Db.GetDatabase(),
		request: payload.Http.Request,
		history: payload.Services.HistoryDb,
	}
}

func (s *RoutesService) List() (*models.SvcRoutes, *logger.HttpError) {

	/* 	models, err := s.loadData()
	   	if err != nil {
	   		return nil, err
	   	} */

	return nil, nil
}

func (s *RoutesService) Refresh() *logger.HttpError {

	logger.Info("refreshing permissions database")
	models, err := s.loadData()
	if err != nil {
		return err
	}

	logger.Info("found %d active routes", len(*models.Docs))
	for _, model := range *models.Docs {

		key := permission_key.MakePermissionKey(model.Method, model.Uri)

		if err := s.payload.Services.PermissionsDb.Set(key, model, 0); err != nil {
			return logger.Error(
				logger.LogStatusInternalServerError,
				nil,
				"failed to set to write top permissions database",
				err,
				nil,
			)
		}
	}

	// TODO: delete unused keys from Redis
	logger.Info("permissions database refreshed")

	return nil
}

func (s *RoutesService) loadData() (*models.SvcRoutes, *logger.HttpError) {

	models := models.SvcRoutes{}

	sql := "SELECT " +
		"a.code as app, " +
		"arm.`method`, " +
		"CONCAT(ar.uri,arm.uri) as uri, " +
		"arm.`open` " +
		"FROM " +
		"applications a " +
		"LEFT JOIN  app_routes ar ON (a.id = ar.application_id)  " +
		"LEFT JOIN app_environments ae ON (ae.id = ar.app_environment_id)  " +
		"LEFT JOIN app_route_methods arm ON (arm.app_route_id = ar.id)  " +
		"WHERE  " +
		"ae.id = ?  " +
		"AND ar.active <> 0  " +
		"AND arm.active <> 0  " +
		"AND a.deleted_by IS NULL " +
		"AND a.deleted_at IS NULL " +
		"AND ae.deleted_by IS NULL " +
		"AND ae.deleted_at IS NULL " +
		"AND ar.deleted_by IS NULL " +
		"AND ar.deleted_at IS NULL " +
		"AND arm.deleted_by IS NULL " +
		"AND arm.deleted_at IS NULL "

	if err := s.db.Raw(
		sql,
		uint(s.payload.Config.GetApplication().EnvironmentID),
	).Scan(&models.Docs).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {

			return nil, logger.Error(
				logger.LogStatusNotFound,
				nil,
				"no documents found",
				nil,
				nil,
			)
		}
		return nil, logger.Error(
			logger.LogStatusInternalServerError,
			nil,
			"failed to query database",
			err,
			nil,
		)
	}

	return &models, nil

}

/* func (s *RoutesService) List(filter string, args ...any) (*models.UserGroups, *logger.HttpError) {

	model := models.UserGroup{}
	models := models.UserGroups{}

	if err := s.db.Model(&model).Where(filter, args).Count(&models.Count).Error; err != nil {

		return nil, logger.Error(
			logger.LogStatusInternalServerError,
			nil,
			"failed to query database",
			err,
			nil,
		)
	}

	if err := s.db.Model(&model).Where(filter, args).Find(&models.Docs).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {

			return nil, logger.Error(
				logger.LogStatusNotFound,
				nil,
				"no documents found",
				nil,
				nil,
			)
		}

		return nil, logger.Error(
			logger.LogStatusInternalServerError,
			nil,
			"failed to query database",
			err,
			nil,
		)
	}

	return &models, nil
} */
