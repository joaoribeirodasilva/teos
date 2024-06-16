package configuration

import (
	"errors"
	"strings"
	"sync"

	"github.com/joaoribeirodasilva/teos/dbtest/database"
	"github.com/joaoribeirodasilva/teos/dbtest/logger"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	ErrValueNotFound      = errors.New("confuration value not found")
	ErrValueDifferentType = errors.New("confuration value exists with a different data type")

	defaultOptions = &ConfigurationOptions{
		Application:                    "",
		Db:                             nil,
		AppConfigurationCollectionName: "app_configurations",
		AppCollectionName:              "app_apps",
	}
)

type ConfigurationApp struct {
	ID  primitive.ObjectID `json:"_id" bson:"_id"`
	Key string             `json:"appKey" bson:"appKey"`
}

type ConfigurationValue struct {
	ID            primitive.ObjectID `json:"_id" bson:"_id"`
	ApplicationID primitive.ObjectID `json:"appAppId" bson:"appAppId"`
	Name          string             `json:"name" bson:"name"`
	Key           string             `json:"key" bson:"key"`
	Type          string             `json:"type" bson:"type"`
	Int           int64              `json:"valueInt" bson:"valueInt"`
	String        string             `json:"valueString" bson:"valueString"`
	Float         float64            `json:"valueFloat" bson:"valueFloat"`
	Bool          bool               `json:"valueBool" bson:"valueBool"`
}

type ConfigurationOptions struct {
	Application                    string
	Db                             *database.Db
	AppConfigurationCollectionName string
	AppCollectionName              string
}

type Configuration struct {
	applicationID primitive.ObjectID
	options       ConfigurationOptions
	values        map[string]ConfigurationValue
	waitgGroup    sync.WaitGroup
}

func New(opts *ConfigurationOptions) *Configuration {

	options := defaultOptions
	if opts != nil {
		options = opts
	}

	c := &Configuration{
		values:  make(map[string]ConfigurationValue),
		options: *options,
	}

	return c
}

func (c *Configuration) GetConfiguration() error {

	logger.Info("reading configuration from database")

	if err := c.GetApp(); err != nil {
		return err
	}

	if err := c.GetAppConfiguraton(); err != nil {
		return err
	}

	logger.Info("confuguration readed from database successfully")

	logger.Debug("configuration now is: ", c.values)

	return nil
}

func (c *Configuration) GetApp() error {

	coll := c.options.Db.GetDatabase().Collection(c.options.AppCollectionName)

	app := &ConfigurationApp{}

	if err := coll.FindOne(c.options.Db.GetContext(), bson.D{{Key: "appKey", Value: c.options.Application}}).Decode(app); err != nil {
		return err
	}

	c.applicationID = app.ID

	logger.Debug("application id is: %s", nil, c.applicationID.Hex())

	return nil
}

func (c *Configuration) GetAppConfiguraton() error {

	coll := c.options.Db.GetDatabase().Collection(c.options.AppConfigurationCollectionName)

	configs := []ConfigurationValue{}

	filter := bson.D{
		{Key: "$and", Value: bson.A{
			bson.D{
				{Key: "$or", Value: bson.A{
					bson.D{{Key: "appAppId", Value: c.applicationID}},
					bson.D{{Key: "appAppId", Value: nil}},
				},
				},
			},
			bson.D{{Key: "deletedBy", Value: nil}},
			bson.D{{Key: "deletedAt", Value: nil}},
		},
		},
	}

	cursor, err := coll.Find(c.options.Db.GetContext(), filter)
	if err != nil {
		return err
	}

	if err := cursor.All(c.options.Db.GetContext(), &configs); err != nil {
		return err
	}

	c.waitgGroup.Add(1)
	for _, conf := range configs {
		key := strings.ToLower(conf.Key)
		conf.Type = strings.ToLower(conf.Type)
		c.values[key] = conf
	}
	c.waitgGroup.Done()

	return nil
}

func (c *Configuration) GetValue(key string) (*ConfigurationValue, error) {

	c.waitgGroup.Wait()

	val, ok := c.values[key]
	if !ok {
		return nil, ErrValueNotFound
	}

	return &val, nil
}

func (c *Configuration) GetInt(key string) (int64, error) {

	c.waitgGroup.Wait()

	val, ok := c.values[key]
	if !ok {
		return 0, ErrValueNotFound
	}

	if val.Type != "int" {
		return 0, ErrValueDifferentType
	}

	return val.Int, nil
}

func (c *Configuration) GetString(key string) (string, error) {

	c.waitgGroup.Wait()

	val, ok := c.values[key]
	if !ok {
		return "", ErrValueNotFound
	}

	if val.Type != "string" {
		return "", ErrValueDifferentType
	}

	return val.String, nil
}

func (c *Configuration) GetFloat(key string) (float64, error) {

	c.waitgGroup.Wait()

	val, ok := c.values[key]
	if !ok {
		return 0, ErrValueNotFound
	}

	if val.Type != "string" {
		return 0, ErrValueDifferentType
	}

	return val.Float, nil

}

func (c *Configuration) GetBool(key string) (bool, error) {

	c.waitgGroup.Wait()

	val, ok := c.values[key]
	if !ok {
		return false, ErrValueNotFound
	}

	if val.Type != "string" {
		return false, ErrValueDifferentType
	}

	return val.Bool, nil
}
