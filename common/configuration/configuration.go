package configuration

import (
	"context"
	"fmt"
	"log/slog"

	app_models "github.com/joaoribeirodasilva/teos/apps/models"
	"github.com/joaoribeirodasilva/teos/common/conf"
	"github.com/joaoribeirodasilva/teos/common/database"
	"github.com/joaoribeirodasilva/teos/common/service_errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ConfigurationValues struct {
	Int    *int
	String *string
	Float  *float64
	Type   string
}

type Configuration struct {
	db     *database.Db
	conf   *conf.Conf
	Config map[string]ConfigurationValues
	ID     primitive.ObjectID
}

func New(db *database.Db, conf *conf.Conf) *Configuration {
	return &Configuration{
		db:     db,
		conf:   conf,
		Config: make(map[string]ConfigurationValues, 0),
	}
}

func (c *Configuration) GetAppId() *service_errors.Error {

	slog.Info("reading service identification from database...")
	application := app_models.AppApp{}
	coll := c.db.Db.Collection("app_apps")
	if err := coll.FindOne(context.TODO(), bson.D{{Key: "appKey", Value: c.conf.Service.Name}}).Decode(&application); err != nil {
		return service_errors.New(0, 0, "CONFIGURATION", "GetAppId", "failed to get service identification from database. ERR: %s", err.Error()).LogError()
	}

	c.ID = application.ID
	slog.Info(fmt.Sprintf("[CONFIGURATION] service identification is: %s\n", c.ID.Hex()))

	return nil

}

func (c *Configuration) Read() error {

	slog.Info("reading service configuration from database...")

	coll := c.db.Db.Collection("app_configurations")
	cursor, err := coll.Find(context.TODO(), bson.D{{Key: "appAppId", Value: c.ID}})
	if err != nil {
		if err != mongo.ErrNoDocuments {
			return service_errors.New(0, 0, "CONFIGURATION", "Read", "failed to get service configuration from database. ERR: %s", err.Error()).LogError()
		}
	}

	configurations := []app_models.AppConfiguration{}
	if err := cursor.All(context.TODO(), &configurations); err != nil {
		return service_errors.New(0, 0, "CONFIGURATION", "Read", "failed to get service configuration from database. ERR: %s", err.Error()).LogError()
	}

	tempConfig := make(map[string]ConfigurationValues, 0)
	for _, config := range configurations {
		tempConfig[config.Key] = ConfigurationValues{
			Type:   config.Type,
			Int:    config.ValueInt,
			String: config.ValueString,
			Float:  config.ValueFloat,
		}
	}
	slog.Info(fmt.Sprintf("service configuration read %d values", len(tempConfig)))

	c.Config = tempConfig

	return nil

}

func (c *Configuration) GetKey(key string) *ConfigurationValues {
	val, ok := c.Config[key]
	if ok {
		return &val
	}
	return nil
}
