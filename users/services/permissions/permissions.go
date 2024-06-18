package permissions

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/joaoribeirodasilva/teos/common/database"
	"github.com/joaoribeirodasilva/teos/common/logger"
	"github.com/joaoribeirodasilva/teos/common/models"
	"github.com/joaoribeirodasilva/teos/common/redisdb"
	"github.com/joaoribeirodasilva/teos/common/requests"
	"github.com/joaoribeirodasilva/teos/common/services"
	"github.com/joaoribeirodasilva/teos/common/structures"
	"github.com/joaoribeirodasilva/teos/common/utils/token"
	"github.com/joaoribeirodasilva/teos/hist/services/histories"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	collectionName = "user_permissions"
)

type PermissionsService struct {
	services      *structures.RequestValues
	database      *database.Db
	db            *mongo.Database
	coll          *mongo.Collection
	user          *token.User
	query         *requests.QueryString
	sessionDb     *redisdb.RedisDB
	permissionsDb *redisdb.RedisDB
	context       context.Context
}

func New(services *structures.RequestValues) *PermissionsService {
	s := &PermissionsService{}
	s.services = services
	s.database = services.Services.Db
	s.db = s.database.GetDatabase()
	s.coll = s.db.Collection(collectionName)
	s.user = services.User
	s.query = &services.Query
	s.sessionDb = services.Services.SessionsDB
	s.permissionsDb = services.Services.PermissionsDB
	s.context = s.database.GetContext()
	return s
}

// List returns a list of user sessions from the collection
func (s *PermissionsService) List(filter bson.D) (*models.UserPermissionsModel, *logger.HttpError) {

	opts := options.FindOptions{}
	if filter == nil {

		filter = *s.query.Filter
		opts = *s.query.Options
	} else {

		opts.SetSkip(0)
		opts.SetLimit(10)
		opts.SetSort(
			bson.D{
				{Key: "created_at", Value: -1},
			},
		)
	}

	count, err := s.coll.CountDocuments(s.context, filter, nil)
	if err != nil {

		return nil, logger.Error(
			logger.LogStatusInternalServerError,
			nil,
			"failed to query database",
			err,
			nil,
		)
	}

	cursor, err := s.coll.Find(s.context, filter, &opts)
	if err != nil {

		if err == mongo.ErrNoDocuments {

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

	docs := &models.UserPermissionsModel{
		Count: count,
	}

	if err := cursor.All(s.context, docs.Docs); err != nil {

		return nil, logger.Error(
			logger.LogStatusInternalServerError,
			nil,
			"failed fetch results",
			err,
			nil,
		)
	}

	return nil, nil
}

// Get returns a single user sessions from the collection
func (s *PermissionsService) Get(filter bson.D, model *models.UserPermissionModel) *logger.HttpError {

	if filter == nil {
		filter = bson.D{{Key: "_id", Value: s.query.ID}}
	}

	if err := s.coll.FindOne(
		s.context,
		filter,
	).Decode(model); err != nil {

		if err == mongo.ErrNoDocuments {
			return logger.Error(logger.LogStatusNotFound, nil, "no documents found", nil, nil)
		}
		return logger.Error(logger.LogStatusInternalServerError, nil, "failed to query database", err, nil)
	}

	return nil
}

// Create creates a new user session document or returns a logger.HttpError in case of error
func (s *PermissionsService) Create(model *models.UserPermissionModel) *logger.HttpError {

	if err := s.Validate(model); err != nil {
		return err
	}

	exists := &models.UserPermissionModel{}
	if err := s.Get(
		bson.D{{Key: "$and", Value: bson.A{
			bson.D{
				{Key: "appRouteId", Value: model.AppRouteID},
			},
			bson.D{{Key: "$or", Value: bson.A{
				bson.D{
					{Key: "userRoleGroupId", Value: model.UserRoleGroupID},
				},
				bson.D{
					{Key: "userRoleGroupId", Value: model.UserUserID},
				},
			},
			}},
		}}},
		exists,
	); err != nil {
		if err.Status != logger.LogStatusNotFound {
			return err
		}
		return nil
	} else {

		historySvc := histories.New(s.services)
		history := &models.HistHistoryModel{
			AppAppID:   s.services.Services.Configuration.GetAppID(),
			Collection: collectionName,
			OriginalID: exists.ID,
			Data:       exists,
		}
		if err := historySvc.Create(history); err != nil {
			return err
		}

		exists.DeletedBy = nil
		exists.DeletedAt = nil
	}

	s.assign(
		exists,
		model,
		services.SVC_OPERATION_CREATE,
	)

	if _, err := s.coll.InsertOne(s.context, model); err != nil {

		return logger.Error(
			logger.LogStatusInternalServerError,
			nil,
			"failed to create document",
			err,
			nil,
		)
	}

	return nil
}

// Create updates a user session document or returns a logger.HttpError in case of error
func (s *PermissionsService) Update(model *models.UserPermissionModel) *logger.HttpError {

	if err := s.Validate(model); err != nil {
		return err
	}

	exists := &models.UserPermissionModel{}

	if err := s.Get(
		bson.D{{Key: "$and", Value: bson.A{
			bson.D{
				{Key: "appRouteId", Value: model.AppRouteID},
			},
			bson.D{{Key: "$or", Value: bson.A{
				bson.D{
					{Key: "userRoleGroupId", Value: model.UserRoleGroupID},
				},
				bson.D{
					{Key: "userRoleGroupId", Value: model.UserUserID},
				},
			},
			}},
		}}},
		exists,
	); err != nil {
		if err != mongo.ErrNoDocuments {
			return err
		}
	}

	if exists.ID != model.ID {
		err := errors.New("duplicate")
		fields := []string{"email"}
		return logger.Error(logger.LogStatusConflict, &fields, fmt.Sprintf("document already exists with id: %s", exists.ID.Hex()), err, nil)
	}

	historySvc := histories.New(s.services)
	history := &models.HistHistoryModel{
		AppAppID:   s.services.Services.Configuration.GetAppID(),
		Collection: collectionName,
		OriginalID: exists.ID,
		Data:       exists,
	}
	if err := historySvc.Create(history); err != nil {
		return err
	}

	s.assign(
		exists,
		model,
		services.SVC_OPERATION_UPDATE,
	)

	if _, err := s.coll.UpdateOne(
		s.context,
		bson.D{{Key: "_id", Value: s.query.ID}},
		bson.D{{Key: "$set", Value: model}},
	); err != nil {
		return logger.Error(logger.LogStatusInternalServerError, nil, "failed to update document", err, nil)
	}

	return nil
}

// Delete deletes a user session document or returns a logger.HttpError in case of error
func (s *PermissionsService) Delete(model *models.UserPermissionModel) *logger.HttpError {

	exists := &models.UserPermissionModel{}

	if err := s.Get(
		bson.D{
			{Key: "_id", Value: s.query.ID},
		},
		exists,
	); err != nil {

		return err
	}

	historySvc := histories.New(s.services)
	history := &models.HistHistoryModel{
		AppAppID:   s.services.Services.Configuration.GetAppID(),
		Collection: collectionName,
		OriginalID: exists.ID,
		Data:       exists,
	}
	if err := historySvc.Create(history); err != nil {
		return err
	}

	s.assign(
		exists,
		model,
		services.SVC_OPERATION_DELETE,
	)

	if _, err := s.coll.UpdateOne(
		s.context,
		bson.D{{Key: "_id", Value: s.query.ID}},
		bson.D{{Key: "$set", Value: model}},
	); err != nil {
		return logger.Error(
			logger.LogStatusInternalServerError,
			nil,
			"failed to update document",
			err,
			nil,
		)
	}

	return nil
}

func (m *PermissionsService) Validate(model *models.UserPermissionModel) *logger.HttpError {

	validate := validator.New()

	if err := validate.Var(model.AppRouteID, "required"); err != nil {
		fields := []string{"appRouteId"}
		return logger.Error(logger.LogStatusBadRequest, &fields, "invalid appRouteId ", err, nil)
	}

	if err := validate.Var(model.UserUserID, "required"); err != nil {
		fields := []string{"userUserId"}
		return logger.Error(logger.LogStatusBadRequest, &fields, "invalid userUserId ", err, nil)
	}

	if err := validate.Var(model.Active, "required"); err != nil {
		fields := []string{"active"}
		return logger.Error(logger.LogStatusBadRequest, &fields, "invalid active ", err, nil)
	}

	// TODO: Validate related
	/*	userAppRoute := NewAppRouteModel(m.ctx)
		 	if appErr := m.FindByID(m.AppRouteID, userAppRoute); appErr != nil {
				return appErr
			}

			userRoleGroup := NewUserRolesGroupModel(m.ctx)
			if appErr := m.FindByID(m.UserRoleGroupID, userRoleGroup); appErr != nil {
				return appErr
			}

			user := NewUserUserModel(m.ctx)
			if appErr := m.FindByID(m.UserUserID, user); appErr != nil {
				return appErr
			}
	*/

	return nil
}

func (s *PermissionsService) assign(to *models.UserPermissionModel, from *models.UserPermissionModel, operation services.Operation) {

	now := time.Now().UTC()

	if operation == services.SVC_OPERATION_CREATE {

		to.ID = primitive.NewObjectID()
		to.CreatedBy = s.user.ID
		to.CreatedAt = now

	} else if operation == services.SVC_OPERATION_DELETE {

		to.DeletedBy = &s.user.ID
		to.DeletedAt = &now

	} else {

		to.AppRouteID = from.AppRouteID
		to.UserRoleGroupID = from.UserRoleGroupID
		to.UserUserID = from.UserUserID
		to.Active = from.Active
	}

	to.UpdatedBy = s.user.ID
	to.UpdatedAt = now
}
