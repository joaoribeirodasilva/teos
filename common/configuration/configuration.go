package configuration

import (
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/joaoribeirodasilva/teos/common/database"
	"github.com/joaoribeirodasilva/teos/common/logger"
)

var (
	ErrValueNotFound      = errors.New("configuration value not found")
	ErrValueDifferentType = errors.New("configuration value exists with a different data type")

	defaultOptions = &ConfigurationOptions{
		Application: "",
		Db:          nil,
	}

	redisInstances = []string{"DB_PERMISSIONS", "DB_SESSIONS", "DB_HISTORY", "DB_LOGS"}
)

type ConfigurationApp struct {
	ID   uint   `json:"_id" gorm:"column:id"`
	Code string `json:"code" gorm:"column:code"`
}

func (ca *ConfigurationApp) TableName() string {
	return "applications"
}

type EnvironmentApp struct {
	ID  uint   `json:"_id" gorm:"column:id"`
	Key string `json:"key" gorm:"column:key"`
}

func (ca *EnvironmentApp) TableName() string {
	return "app_environments"
}

type ConfigurationValue struct {
	ID               uint       `json:"id" gorm:"column:id;type:uint;primaryKey"`
	AppEnvironmentID uint       `json:"appEnvironmentId" gorm:"column:app_environment_id;type:uint;not null;"`
	ApplicationID    uint       `json:"applicationId" gorm:"column:application_ịId;type:uint;not null;"`
	ConfigurationKey string     `json:"configurationKey" gorm:"column:configuration_key;type:string;size:255;not null;"`
	ValString        *string    `json:"valString" gorm:"column:val_string;type:string;size:65536;"`
	ValInt           *int64     `json:"valInt" gorm:"column:val_int;type:int64;"`
	ValDouble        *float64   `json:"valDouble" gorm:"column:val_double;type:float64;"`
	ValBool          *int       `json:"valBoolean" gorm:"column:val_boolean;type:int;size(1);"`
	ValDate          *time.Time `json:"valDate" gorm:"column:val_date;type:DATE;"`
	ValTime          *time.Time `json:"valTime" gorm:"column:val_time;type:TIME;"`
	ValDateTime      *time.Time `json:"valDateTime" gorm:"column:val_datetime;type:TIMESTAMP;"`
	Type             string     `json:"type" gorm:"column:type;type:string;"`
}

func (ca *ConfigurationValue) TableName() string {
	return "app_configurations"
}

type ConfigurationOptions struct {
	Application string
	Environment string
	Db          *database.Db
}

type ConfigurationRedis struct {
	Addr     string
	Port     int
	Db       int
	Username string
	Password string
}

type Configuration struct {
	Secret        string
	applicationID uint
	environmentID uint
	options       ConfigurationOptions
	values        map[string]ConfigurationValue
	redis         map[string]ConfigurationRedis
	waitgGroup    sync.WaitGroup
}

func New(opts *ConfigurationOptions) *Configuration {

	options := defaultOptions
	if opts != nil {
		options = opts
	}

	c := &Configuration{
		values:  make(map[string]ConfigurationValue),
		redis:   make(map[string]ConfigurationRedis),
		options: *options,
	}

	return c
}

func (c *Configuration) GetConfiguration() error {

	logger.Info("reading configuration from database")

	if err := c.GetApp(); err != nil {
		return err
	}

	if err := c.GetEnvironment(); err != nil {
		return err
	}

	if err := c.GetAppConfiguration(); err != nil {
		return err
	}

	logger.Info("configuration reded from database successfully")

	logger.Debug("configuration now is: ", c.values)

	return nil
}

func (c *Configuration) GetApp() error {

	app := ConfigurationApp{}

	if err := c.options.Db.GetDatabase().Model(&app).Where("code=?", c.options.Application).First(&app).Error; err != nil {
		return err
	}

	c.applicationID = app.ID

	logger.Debug("application id is: %d", nil, c.applicationID)

	return nil
}

func (c *Configuration) GetEnvironment() error {

	env := EnvironmentApp{}
	if err := c.options.Db.GetDatabase().Model(&env).Where("`key`=?", c.options.Environment).First(&env).Error; err != nil {
		return err
	}

	c.environmentID = env.ID

	logger.Debug("environment id is: %d", nil, c.environmentID)

	return nil
}

func (c *Configuration) GetAppConfiguration() error {

	configs := []ConfigurationValue{}

	if err := c.options.Db.GetDatabase().Model(&configs).Where("(application_id = ? OR application_id=1) AND app_environment_id = ?", c.applicationID, c.environmentID).Find(&configs).Error; err != nil {
		return err
	}

	c.waitgGroup.Add(1)
	for _, conf := range configs {
		key := strings.ToUpper(conf.ConfigurationKey)
		conf.Type = strings.ToLower(conf.Type)
		c.values[key] = conf
		if key == "SECRET_KEY" {
			c.Secret = *conf.ValString
		}
	}

	for _, instance := range redisInstances {

		addr := ""
		port := 0
		db := 0
		user := ""
		pwd := ""

		valString, ok := c.values[instance+"_ADDR"]
		if ok && valString.ValString != nil {
			addr = *valString.ValString
		}

		valInt, ok := c.values[instance+"_PORT"]
		if ok && valInt.ValInt != nil {
			port = int(*valInt.ValInt)
		}

		valInt, ok = c.values[instance+"_DB"]
		if ok && valInt.ValInt != nil {
			db = int(*valInt.ValInt)
		}

		valString, ok = c.values[instance+"_USERNAME"]
		if ok && valString.ValString != nil {
			user = *valString.ValString
		}

		valString, ok = c.values[instance+"_PASSWORD"]
		if ok && valString.ValString != nil {
			pwd = *valString.ValString
		}

		c.redis[instance] = ConfigurationRedis{
			Addr:     addr,
			Port:     int(port),
			Db:       int(db),
			Username: user,
			Password: pwd,
		}

	}

	c.waitgGroup.Done()

	return nil
}

func (c *Configuration) GetAppID() uint {

	return c.applicationID
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
	if !ok || val.ValInt == nil {
		return 0, ErrValueNotFound
	}

	if val.Type != "int" {
		return 0, ErrValueDifferentType
	}

	return *val.ValInt, nil
}

func (c *Configuration) GetString(key string) (string, error) {

	c.waitgGroup.Wait()

	val, ok := c.values[key]
	if !ok || val.ValString == nil {
		return "", ErrValueNotFound
	}

	if val.Type != "string" {
		return "", ErrValueDifferentType
	}

	return *val.ValString, nil
}

func (c *Configuration) GetFloat(key string) (float64, error) {

	c.waitgGroup.Wait()

	val, ok := c.values[key]
	if !ok || val.ValDouble == nil {
		return 0, ErrValueNotFound
	}

	if val.Type != "double" {
		return 0, ErrValueDifferentType
	}

	return *val.ValDouble, nil

}

func (c *Configuration) GetBool(key string) (bool, error) {

	c.waitgGroup.Wait()

	val, ok := c.values[key]
	if !ok || val.ValBool == nil {
		return false, ErrValueNotFound
	}

	if val.Type != "string" {
		return false, ErrValueDifferentType
	}

	ret := false
	if *val.ValBool > 0 {
		ret = true
	}

	return ret, nil
}

func (c *Configuration) GetRedisConf(key string) (*ConfigurationRedis, error) {
	val, ok := c.redis[key]
	if !ok {
		return nil, errors.New("invalid redis database key")
	}

	if val.Addr == "" {
		return nil, fmt.Errorf("redis database %s has no address", key)
	}

	if val.Port == 0 {
		return nil, fmt.Errorf("redis database %s has no port", key)
	}

	return &val, nil
}
