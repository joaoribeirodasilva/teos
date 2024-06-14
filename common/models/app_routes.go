package models

import (
	"net/http"
	"slices"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/joaoribeirodasilva/teos/common/service_errors"
	"github.com/joaoribeirodasilva/teos/common/service_log"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	validMethods = []string{"GET", "PUT", "PATCH", "POST", "DELETE"}
)

const (
	collectionAppRoute = "app_routes"
	locationAppRoute   = "COMMON::MODELS::AppRoute"
)

type AppRouteModel struct {
	BaseModel
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

type AppRouteModels struct {
	BaseModels
	Count int64            `json:"count"`
	Rows  *[]AppRouteModel `json:"rows"`
}

func NewAppRouteModel(c *gin.Context) *AppRouteModel {
	m := &AppRouteModel{}
	m.Init(c, locationAppRoute, collectionAppRoute)
	return m
}

func NewAppRoutesModel(c *gin.Context) *AppRouteModel {
	m := &AppRouteModel{}
	m.Init(c, locationAppRoute, collectionAppRoute)
	return m
}

func (m *AppRouteModel) FillMeta(create bool, delete bool) {

	now := time.Now().UTC()

	if create {
		m.ID = primitive.NewObjectID()
		m.CreatedBy = m.GetValues().Variables.User.ID
		m.CreatedAt = now
	} else if delete {
		m.DeletedBy = &m.GetValues().Variables.User.ID
		m.DeletedAt = &now
	}

	m.UpdatedBy = m.GetValues().Variables.User.ID
	m.UpdatedAt = now
}

func (m *AppRouteModel) Bind() *service_errors.Error {
	return m.BaseModel.Bind(m, m.ctx)
}

func (m *AppRouteModel) Validate() *service_errors.Error {

	validate := validator.New()

	appRoutesBlockID := NewAppRoutesBlockModels(m.ctx)
	if appErr := m.FindByID(m.AppRoutesBlockID, appRoutesBlockID); appErr != nil {
		return appErr
	}

	if err := validate.Var(m.Name, "required,gte=1"); err != nil {
		return service_log.Error(0, http.StatusBadRequest, "MODELS::AppRouteModel", "name", "invalid name")
	}

	if err := validate.Var(m.Method, "required,gte=1"); err != nil {
		return service_log.Error(0, http.StatusBadRequest, "MODELS::AppRouteModel", "method", "invalid method")
	}

	if !slices.Contains(validMethods, m.Method) {
		return service_log.Error(0, http.StatusBadRequest, "MODELS::AppRouteModel", "method", "invalid method")
	}

	if err := validate.Var(m.Route, "required"); err != nil {
		return service_log.Error(0, http.StatusBadRequest, "MODELS::AppRouteModel", "route", "invalid route")
	}

	if err := validate.Var(m.Open, "required"); err != nil {
		return service_log.Error(0, http.StatusBadRequest, "MODELS::AppRouteModel", "open", "invalid open")
	}

	if err := validate.Var(m.Active, "required"); err != nil {
		return service_log.Error(0, http.StatusBadRequest, "MODELS::AppRouteModel", "active", "invalid active")
	}

	return nil
}
