package histories

import (
	"context"
	"time"

	"github.com/joaoribeirodasilva/teos/common/database"
	"github.com/joaoribeirodasilva/teos/common/logger"
	"github.com/joaoribeirodasilva/teos/common/models"
	"github.com/joaoribeirodasilva/teos/common/redisdb"
	"github.com/joaoribeirodasilva/teos/common/requests"
	"github.com/joaoribeirodasilva/teos/common/services"
	"github.com/joaoribeirodasilva/teos/common/structures"
	"github.com/joaoribeirodasilva/teos/common/utils/token"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	collectionName = "user_sessions"
)

type HistoriesService struct {
	database      *database.Db
	db            *mongo.Database
	coll          *mongo.Collection
	user          *token.User
	query         *requests.QueryString
	sessionDb     *redisdb.RedisDB
	permissionsDb *redisdb.RedisDB
	context       context.Context
}

func New(services *structures.RequestValues) *HistoriesService {
	s := &HistoriesService{}
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
func (s *HistoriesService) List(filter bson.D) (*models.HistHistoriesModel, *logger.HttpError) {

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

	docs := &models.HistHistoriesModel{
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
func (s *HistoriesService) Get(filter bson.D, model *models.HistHistoryModel) *logger.HttpError {

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
func (s *HistoriesService) Create(model *models.HistHistoryModel) *logger.HttpError {

	if err := s.Validate(model); err != nil {
		return err
	}

	s.assign(
		nil,
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
func (s *HistoriesService) Update(model *models.HistHistoryModel) *logger.HttpError {

	if err := s.Validate(model); err != nil {
		return err
	}

	exists := &models.HistHistoryModel{}

	if err := s.Get(
		bson.D{
			{Key: "_id", Value: s.query.ID},
		},
		exists,
	); err != nil {

		return err
	}

	// TODO: History

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
func (s *HistoriesService) Delete(model *models.HistHistoryModel) *logger.HttpError {

	exists := &models.HistHistoryModel{}

	if err := s.Get(
		bson.D{
			{Key: "_id", Value: s.query.ID},
		},
		exists,
	); err != nil {

		return err
	}

	// TODO: History

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

func (m *HistoriesService) Validate(model *models.HistHistoryModel) *logger.HttpError {

	// TODO: Validate related

	/* 	user := NewUserUserModel(m.ctx)
	   	if appErr := m.FindByID(m.UserUserID, user); appErr != nil {
	   		return appErr
	   	} */

	return nil
}

func (s *HistoriesService) assign(to *models.HistHistoryModel, from *models.HistHistoryModel, operation services.Operation) {

	now := time.Now().UTC()

	if operation == services.SVC_OPERATION_CREATE {

		to.ID = primitive.NewObjectID()
		to.CreatedBy = s.user.ID
		to.CreatedAt = now

	} else if operation == services.SVC_OPERATION_DELETE {

		to.DeletedBy = &s.user.ID
		to.DeletedAt = &now

	} else {

		to.AppAppID = from.AppAppID
		to.Collection = from.Collection
		to.OriginalID = from.OriginalID
		to.Data = from.Data
	}

	to.UpdatedBy = s.user.ID
	to.UpdatedAt = now
}
