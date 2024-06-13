package configuration

import (
	"context"
	"fmt"
	"log/slog"

	app_models "github.com/joaoribeirodasilva/teos/apps/models"
	"github.com/joaoribeirodasilva/teos/common/conf"
	"github.com/joaoribeirodasilva/teos/common/database"
	"github.com/joaoribeirodasilva/teos/common/service_errors"
	"github.com/joaoribeirodasilva/teos/common/service_log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ConfigurationValues struct {
	Int    *int
	String *string
	Float  *float64
	Bool   *bool
	Type   string
}

type ConfigRedis struct {
	Addresses string
	Username  string
	Password  string
	Db        int
}

type Configuration struct {
	db                *database.Db
	conf              *conf.Conf
	Config            map[string]ConfigurationValues
	ID                primitive.ObjectID
	DbHistory         ConfigRedis
	DbLog             ConfigRedis
	DbPermissions     ConfigRedis
	DbSessions        ConfigRedis
	appConfAppKeys    map[string]*ConfigApp
	appConfigAppAuths map[string]*ConfigApp
	appConfigAppIDs   map[primitive.ObjectID]*ConfigApp
}

type ConfigApp struct {
	ID      primitive.ObjectID
	Name    string
	Key     string
	AuthKey string
	DnsName string
	Port    int
}

func New(db *database.Db, conf *conf.Conf) *Configuration {
	return &Configuration{
		db:     db,
		conf:   conf,
		Config: make(map[string]ConfigurationValues, 0),
	}
}

func (c *Configuration) GetAppId() *service_errors.Error {

	slog.Info("[COMMON::CONFIGURATION::GetAppId] reading service identification from database...")
	application := app_models.AppApp{}
	coll := c.db.Db.Collection("app_apps")
	if err := coll.FindOne(context.TODO(), bson.D{{Key: "appKey", Value: c.conf.Service.Name}}).Decode(&application); err != nil {
		return service_log.Error(0, 0, "COMMON::CONFIGURATION::GetAppId", "", "failed to get service identification from database. ERR: %s", err.Error())
	}

	c.ID = application.ID
	slog.Info(fmt.Sprintf("[COMMON::CONFIGURATION::GetAppId] service identification is: %s\n", c.ID.Hex()))

	return nil

}

func (c *Configuration) Read() error {

	slog.Info("[COMMON::CONFIGURATION::Read] reading service configuration from database...")

	coll := c.db.Db.Collection("app_configurations")
	cursor, err := coll.Find(context.TODO(), bson.D{{Key: "appAppId", Value: c.ID}})
	if err != nil {
		if err != mongo.ErrNoDocuments {
			return service_log.Error(0, 0, "COMMON::CONFIGURATION::Read", "", "failed to get service configuration from database. ERR: %s", err.Error())
		}
	}

	configurations := []app_models.AppConfiguration{}
	if err := cursor.All(context.TODO(), &configurations); err != nil {
		return service_log.Error(0, 0, "COMMON::CONFIGURATION::Read", "", "failed to get service configuration from database. ERR: %s", err.Error())
	}

	tempConfig := make(map[string]ConfigurationValues, 0)
	for _, config := range configurations {
		tempConfig[config.Key] = ConfigurationValues{
			Type:   config.Type,
			Int:    config.ValueInt,
			String: config.ValueString,
			Float:  config.ValueFloat,
			Bool:   config.ValueBool,
		}
	}
	slog.Info(fmt.Sprintf("[COMMON::CONFIGURATION::Read] service configuration read %d values", len(tempConfig)))

	c.Config = tempConfig

	if appErr := c.loadAppKeys(); appErr != nil {
		return appErr
	}

	if appErr := c.loadGlobalConfig(); appErr != nil {
		return appErr
	}

	return nil

}

func (c *Configuration) loadGlobalConfig() *service_errors.Error {

	coll := c.db.Db.Collection("app_configurations")
	cursor, err := coll.Find(context.TODO(), bson.D{{Key: "appAppId", Value: nil}})
	if err != nil {
		if err != mongo.ErrNoDocuments {
			return nil
		}
		return service_log.Error(0, 0, "COMMON::CONFIGURATION::Read", "", "failed to query database. ERR: %s", err.Error())
	}

	configs := []app_models.AppConfiguration{}
	if err = cursor.All(context.TODO(), &configs); err != nil {
		return service_log.Error(0, 0, "COMMON::CONFIGURATION::Read", "", "failed to fetch result. ERR: %s", err.Error())
	}

	for _, config := range configs {
		switch config.Key {
		case "REDIS_DB_HIST_ADDRESS":
			if config.ValueString == nil {
				continue
			}
			c.DbHistory.Addresses = *config.ValueString
		case "REDIS_DB_HIST_USERNAME":
			if config.ValueString == nil {
				continue
			}
			c.DbHistory.Username = *config.ValueString
		case "REDIS_DB_HIST_PASSWORD":
			if config.ValueString == nil {
				continue
			}
			c.DbHistory.Password = *config.ValueString
		case "REDIS_DB_HIST_DATABASE":
			if config.ValueInt == nil {
				continue
			}
			c.DbHistory.Db = *config.ValueInt
		case "REDIS_DB_LOG_ADDRESS":
			if config.ValueString == nil {
				continue
			}
			c.DbLog.Addresses = *config.ValueString
		case "REDIS_DB_LOG_USERNAME":
			if config.ValueString == nil {
				continue
			}
			c.DbLog.Username = *config.ValueString
		case "REDIS_DB_LOG_PASSWORD":
			if config.ValueString == nil {
				continue
			}
			c.DbLog.Password = *config.ValueString
		case "REDIS_DB_LOG_DATABASE":
			if config.ValueInt == nil {
				continue
			}
			c.DbLog.Db = *config.ValueInt
		case "REDIS_DB_SESSIONS_ADDRESS":
			if config.ValueString == nil {
				continue
			}
			c.DbSessions.Addresses = *config.ValueString
		case "REDIS_DB_SESSIONS_USERNAME":
			if config.ValueString == nil {
				continue
			}
			c.DbSessions.Username = *config.ValueString
		case "REDIS_DB_SESSIONS_PASSWORD":
			if config.ValueString == nil {
				continue
			}
			c.DbSessions.Password = *config.ValueString
		case "REDIS_DB_SESSIONS_DATABASE":
			if config.ValueInt == nil {
				continue
			}
			c.DbSessions.Db = *config.ValueInt
		case "REDIS_DB_PERMISSIONS_ADDRESS":
			if config.ValueString == nil {
				continue
			}
			c.DbPermissions.Addresses = *config.ValueString
		case "REDIS_DB_PERMISSIONS_USERNAME":
			if config.ValueString == nil {
				continue
			}
			c.DbPermissions.Username = *config.ValueString
		case "REDIS_DB_PERMISSIONS_PASSWORD":
			if config.ValueString == nil {
				continue
			}
			c.DbPermissions.Password = *config.ValueString
		case "REDIS_DB_PERMISSIONS_DATABASE":
			if config.ValueInt == nil {
				continue
			}
			c.DbPermissions.Db = *config.ValueInt
		}
	}

	return nil
}

func (c *Configuration) loadAppKeys() *service_errors.Error {

	coll := c.db.Db.Collection("app_apps")
	cursor, err := coll.Find(context.TODO(), bson.D{{}})
	if err != nil {
		if err != mongo.ErrNoDocuments {
			return nil
		}
		return service_log.Error(0, 0, "COMMON::CONFIGURATION::loadAppKeys", "", "failed to query database. ERR: %s", err.Error())
	}

	applications := []app_models.AppApp{}
	if err := cursor.All(context.TODO(), &applications); err != nil {
		return service_log.Error(0, 0, "COMMON::CONFIGURATION::loadAppKeys", "", "failed to fetch result. ERR: %s", err.Error())
	}

	c.appConfigAppIDs = make(map[primitive.ObjectID]*ConfigApp)
	c.appConfAppKeys = make(map[string]*ConfigApp)
	for _, application := range applications {
		app := &ConfigApp{
			ID:      application.ID,
			Name:    application.Name,
			Key:     application.AppKey,
			AuthKey: "",
			DnsName: "",
			Port:    0,
		}

		c.appConfigAppIDs[application.ID] = app
		c.appConfAppKeys[application.AppKey] = app
	}

	coll = c.db.Db.Collection("app_configurations")
	for id, application := range c.appConfigAppIDs {
		configs := []app_models.AppConfiguration{}
		cursor, err := coll.Find(context.TODO(), bson.D{{Key: "_id", Value: id}})
		if err != nil {
			if err != mongo.ErrNoDocuments {
				return nil
			}
			return service_log.Error(0, 0, "COMMON::CONFIGURATION::loadAppKeys", "", "failed to query database. ERR: %s", err.Error())
		}

		if err = cursor.All(context.TODO(), &configs); err != nil {
			return service_log.Error(0, 0, "COMMON::CONFIGURATION::loadAppKeys", "", "failed to fetch result. ERR: %s", err.Error())
		}

		for _, conf := range configs {
			switch conf.Key {
			case "APP_AUTH_KEY":
				if conf.ValueString == nil {
					continue
				}
				c.appConfigAppIDs[application.ID].AuthKey = *conf.ValueString
				c.appConfAppKeys[application.Key].AuthKey = *conf.ValueString
			case "NET_DNS":
				if conf.ValueString == nil {
					continue
				}
				c.appConfigAppIDs[application.ID].DnsName = *conf.ValueString
				c.appConfAppKeys[application.Key].DnsName = *conf.ValueString
			case "NET_PORT":
				if conf.ValueInt == nil {
					continue
				}
				c.appConfigAppIDs[application.ID].Port = *conf.ValueInt
				c.appConfAppKeys[application.Key].Port = *conf.ValueInt
			}
		}
	}

	c.appConfigAppAuths = make(map[string]*ConfigApp)
	for id, application := range c.appConfigAppIDs {
		c.appConfigAppAuths[application.AuthKey] = c.appConfigAppIDs[id]
	}

	return nil
}

func (c *Configuration) GetKey(key string) *ConfigurationValues {
	val, ok := c.Config[key]
	if ok {
		return &val
	}
	return nil
}

func (c *Configuration) GetAppByAuth(key string) *ConfigApp {

	val, ok := c.appConfigAppAuths[key]
	if !ok {
		return nil
	}
	return val
}

func (c *Configuration) GetAppById(key primitive.ObjectID) *ConfigApp {

	val, ok := c.appConfigAppIDs[key]
	if !ok {
		return nil
	}
	return val
}

func (c *Configuration) GetAppKey(key string) *ConfigApp {

	val, ok := c.appConfAppKeys[key]
	if !ok {
		return nil
	}
	return val
}
