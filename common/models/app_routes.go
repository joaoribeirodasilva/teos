package models

import (
	"errors"
	"slices"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/joaoribeirodasilva/teos/common/logger"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	validMethods = []string{"GET", "PUT", "PATCH", "POST", "DELETE"}
)

const (
	collectionAppRoute = "app_routes"
)

type AppRouteModel struct {
	ID               primitive.ObjectID  `json:"_id" bson:"_id"`
	AppRoutesBlockID primitive.ObjectID  `json:"appRoutesBlockId" bson:"appRoutesBlockId"`
	AppRoutesBlock   AppRoutesBlockModel `json:"appRoutesBlock,omitempty" bson:"-"`
	Name             string              `json:"name" bson:"name"`
	Description      *string             `json:"description" bson:"description"`
	Method           string              `json:"method" bson:"method"`
	Route            string              `json:"route" bson:"route"`
	Open             bool                `json:"open" bson:"open"`
	Active           bool                `json:"active" bson:"active"`
	CreatedBy        primitive.ObjectID  `json:"createdBy" bson:"createdBy"`
	CreatedAt        time.Time           `json:"createdAt" bson:"createdAt"`
	UpdatedBy        primitive.ObjectID  `json:"updatedBy" bson:"updatedBy"`
	UpdatedAt        time.Time           `json:"updatedAt" bson:"updatedAt"`
	DeletedBy        *primitive.ObjectID `json:"deletedBy" bson:"deletedBy"`
	DeletedAt        *time.Time          `json:"deletedAt" bson:"deletedAt"`
}

func (m *AppRouteModel) GetCollectionName() string {
	return collectionAppRoute
}

func (m *AppRouteModel) Validate() *logger.HttpError {

	validate := validator.New()

	// TODO: Validate related
	/* 	appRoutesBlockID := NewAppRoutesBlockModels(m.ctx)
	   	if appErr := m.FindByID(m.AppRoutesBlockID, appRoutesBlockID); appErr != nil {
	   		return appErr
	   	} */

	if err := validate.Var(m.Name, "required,gte=1"); err != nil {
		fields := []string{"name"}
		return logger.Error(logger.LogStatusBadRequest, &fields, "invalid name ", err, nil)
	}

	if err := validate.Var(m.Method, "required,gte=1"); err != nil {
		fields := []string{"method"}
		return logger.Error(logger.LogStatusBadRequest, &fields, "invalid method ", err, nil)
	}

	if !slices.Contains(validMethods, m.Method) {
		fields := []string{"method"}
		err := errors.New("invalid method")
		return logger.Error(logger.LogStatusBadRequest, &fields, "invalid method ", err, nil)
	}

	if err := validate.Var(m.Route, "required"); err != nil {
		fields := []string{"route"}
		return logger.Error(logger.LogStatusBadRequest, &fields, "invalid route ", err, nil)
	}

	if err := validate.Var(m.Open, "required"); err != nil {
		fields := []string{"open"}
		return logger.Error(logger.LogStatusBadRequest, &fields, "invalid open ", err, nil)
	}

	if err := validate.Var(m.Active, "required"); err != nil {
		fields := []string{"active"}
		return logger.Error(logger.LogStatusBadRequest, &fields, "invalid active ", err, nil)
	}

	return nil
}
