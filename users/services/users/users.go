package services

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
	collectionName = "user_users"
)

type UsersService struct {
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

func New(services *structures.RequestValues) *UsersService {
	s := &UsersService{}
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

// List returns a list of users from the collection
func (s *UsersService) List(filter bson.D) (*models.UserUsersModel, *logger.HttpError) {

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

	docs := &models.UserSessionsModel{
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

// Get returns a single user from the collection
func (s *UsersService) Get(filter bson.D, model *models.UserUserModel) *logger.HttpError {

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

// Create creates a new user document or returns a logger.HttpError in case of error
func (s *UsersService) Create(model *models.UserUserModel) *logger.HttpError {

	if err := s.Validate(model); err != nil {
		return err
	}

	exists := &models.UserUserModel{}
	if err := s.Get(
		bson.D{
			{Key: "email", Value: model.Email},
		},
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

// Create updates a user document or returns a logger.HttpError in case of error
func (s *UsersService) Update(model *models.UserUserModel) *logger.HttpError {

	if err := s.Validate(model); err != nil {
		return err
	}

	exists := &models.UserUserModel{}

	if err := s.Get(
		bson.D{
			{Key: "email", Value: model.Email},
		},
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

// Delete deletes a user document or returns a logger.HttpError in case of error
func (s *UsersService) Delete(model *models.UserUserModel) *logger.HttpError {

	exists := &models.UserUserModel{}

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

func (m *UsersService) Validate(model *models.UserUserModel) *logger.HttpError {

	validate := validator.New()
	if err := validate.Var(model.FirstName, "required"); err != nil {
		fields := []string{"firstName"}
		return logger.Error(logger.LogStatusBadRequest, &fields, "invalid firstName ", err, nil)
	}

	if err := validate.Var(model.Surname, "required,gte=1"); err != nil {
		fields := []string{"surname"}
		return logger.Error(logger.LogStatusBadRequest, &fields, "invalid surname ", err, nil)
	}

	if err := validate.Var(model.Email, "required,email"); err != nil {
		fields := []string{"email"}
		return logger.Error(logger.LogStatusBadRequest, &fields, "invalid email ", err, nil)
	}

	if err := validate.Var(model.Password, "required,gte=6"); err != nil {
		fields := []string{"password"}
		return logger.Error(logger.LogStatusBadRequest, &fields, "invalid password ", err, nil)
	}

	if err := validate.Var(model.Terms, "required"); err != nil {
		fields := []string{"terms"}
		return logger.Error(logger.LogStatusBadRequest, &fields, "invalid terms ", err, nil)
	}

	return nil
}

func (s *UsersService) assign(to *models.UserUserModel, from *models.UserUserModel, operation services.Operation) {

	now := time.Now().UTC()

	if operation == services.SVC_OPERATION_CREATE {

		to.ID = primitive.NewObjectID()
		to.CreatedBy = s.user.ID
		to.CreatedAt = now

	} else if operation == services.SVC_OPERATION_DELETE {

		to.DeletedBy = &s.user.ID
		to.DeletedAt = &now

	} else {

		to.FirstName = from.FirstName
		to.Surname = from.Surname
		to.Email = from.Email
		to.Password = from.Password
		to.Terms = from.Terms
		to.AvatarUrl = from.AvatarUrl
		to.EmailVerified = from.EmailVerified
		to.Active = from.Active
	}

	to.UpdatedBy = s.user.ID
	to.UpdatedAt = now
}
