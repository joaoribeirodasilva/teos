package configuration

import (
	"errors"
	"fmt"
	"log/slog"

	app_models "github.com/joaoribeirodasilva/teos/apps/models"
	"github.com/joaoribeirodasilva/teos/common/conf"
	"github.com/joaoribeirodasilva/teos/common/database"
	"gorm.io/gorm"
)

type ConfigurationValues struct {
	Int   *int
	Char  *string
	Float *float64
}

type Configuration struct {
	db     *database.Db
	conf   *conf.Conf
	Config map[string]ConfigurationValues
	ID     uint
}

func New(db *database.Db, conf *conf.Conf) *Configuration {
	return &Configuration{
		db:     db,
		conf:   conf,
		Config: make(map[string]ConfigurationValues, 0),
	}
}

func (c *Configuration) GetAppId() error {

	slog.Info("reading service identification from database...")
	application := app_models.AppApp{}
	if err := c.db.Conn.Where("app_key = ?", c.conf.Service.Name).First(&application).Error; err != nil {
		slog.Error("[CONFIGURATION] failed to get service information from database")
		return err
	}

	c.ID = application.ID
	slog.Info(fmt.Sprintf("service identification is: %d\n", c.ID))

	return nil

}

func (c *Configuration) Read() error {

	configurations := []app_models.AppConfiguration{}

	slog.Info("reading service configuration from database...")

	if err := c.db.Conn.Where("app_app_id = ?", c.ID).Find(&configurations).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			slog.Error("[CONFIGURATION] failed to get service configuration from database")
			return err
		}
	}

	tempConfig := make(map[string]ConfigurationValues, 0)
	for _, config := range configurations {
		tempConfig[config.Key] = ConfigurationValues{
			Int:   config.ValueInt,
			Char:  config.ValueChar,
			Float: config.ValueFloat,
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
