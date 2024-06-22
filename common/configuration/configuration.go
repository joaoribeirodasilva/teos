package configuration

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"strings"

	"github.com/joaoribeirodasilva/teos/common/database"
	"github.com/joho/godotenv"
)

const (
	defaultBindIp          = "127.0.0.1"
	defaultLogLevel        = 3
	defaultDatabaseOptions = ""
)

var (
	ErrNoProtocol         = errors.New("no DB_DSN and no DB_PROTOCOL defined")
	ErrNoHosts            = errors.New("no DB_DSN and no DB_HOSTS defined")
	ErrNoPort             = errors.New("no DB_DSN and no DB_PORT defined")
	ErrNoDatabaseName     = errors.New("no DB_DSN and no DB_DATABASE defined")
	ErrJSONFailed         = errors.New("failed to generate JSON string for the configuration object")
	ErrValueNotFound      = errors.New("configuration value not found")
	ErrValueDifferentType = errors.New("configuration value exists with a different data type")

	//redisInstances = []string{"DB_PERMISSIONS", "DB_SESSIONS", "DB_HISTORY", "DB_LOGS"}
)

type ConfigKeys struct {
	Name    string
	ValType string
	Val     interface{}
}

type Config struct {
	application ConfigApplication
	services    ConfigServices
	log         ConfigLog
	db          *database.Db
	keys        map[string]ConfigKeys
}

func New(applicationCode string) *Config {

	c := &Config{}
	c.application.Code = applicationCode
	c.keys = make(map[string]ConfigKeys)

	return c
}

func (c *Config) GetEnv() error {

	var err error

	slog.Info("reading environment variables")

	if err = godotenv.Load(); err != nil {
		slog.Warn("failed to find .env file, collecting env from real environment")
	}

	c.getEnvLogs()
	c.getApplicationEnv()
	if err := c.getEnvDatabase(); err != nil {
		return err
	}

	slog.Info("environment variables reded successfully")

	return nil
}

func (c *Config) GetConfiguration(db *database.Db) error {

	c.db = db

	if err := c.getApplicationID(); err != nil {
		return err
	}

	if err := c.getConfiguration(); err != nil {
		return err
	}

	if err := c.getServicesConfigurations(); err != nil {
		return err
	}

	return nil
}

func (c *Config) GetApplication() *ConfigApplication {
	return &c.application
}

func (c *Config) GetServices() *ConfigServices {
	return &c.services
}

func (c *Config) GetLog() *ConfigLog {
	return &c.log
}

func (c *Config) GetDatabase() *ConfigDatabase {
	return &c.services.Db
}

func (c *Config) GetPermissions() *ConfigMemDb {
	return &c.services.MemDbs.Permissions
}

func (c *Config) GetSessions() *ConfigMemDb {
	return &c.services.MemDbs.Sessions
}

func (c *Config) GetHistory() *ConfigMemDb {
	return &c.services.MemDbs.History
}

func (c *Config) GetLogs() *ConfigMemDb {
	return &c.services.MemDbs.Logs
}

func (c *Config) getApplicationEnv() {

	c.application.Code = os.Getenv("SERVICE_NAME")
	c.application.EnvironmentID = EnvironmentKeys[os.Getenv("ENVIRONMENT")]
	if c.application.EnvironmentID == 0 {
		c.application.EnvironmentID = defaultEnvironmentID
	}
}

func (c *Config) getApplicationID() error {

	app := ConfigurationApp{}

	if err := c.db.GetDatabase().Model(&app).Where("code=?", c.application.Code).First(&app).Error; err != nil {
		return err
	}

	c.application.ID = app.ID

	return nil
}

func (c *Config) getConfiguration() error {

	configs := []ConfigurationValue{}

	if err := c.db.GetDatabase().Model(&configs).Where("(application_id = ? OR application_id=1) AND app_environment_id = ?", c.application.ID, c.application.EnvironmentID).Find(&configs).Error; err != nil {
		return err
	}

	var val interface{}
	for _, config := range configs {

		valType := strings.ToLower(config.Type)
		key := strings.ToUpper(config.ConfigurationKey)

		switch valType {
		case "string":
			val = config.ValString
		case "int":
			val = config.ValInt
		case "double":
			val = config.ValDouble
		case "bool":
			val = config.ValBool
		case "date":
			val = config.ValDate
		case "time":
			val = config.ValTime
		case "datetime":
			val = config.ValDateTime
		}
		c.keys[key] = ConfigKeys{
			Name:    key,
			ValType: valType,
			Val:     val,
		}
	}

	return nil
}

func (c *Config) getServicesConfigurations() error {

	// Http
	conf, ok := c.keys["HTTP_BIND_ADDR"]
	if ok {
		if conf.ValType != "string" {
			return fmt.Errorf("invalid value type %s for configuration key %s", conf.ValType, "HTTP_BIND_ADDR (must be a string)")
		}
	}
	c.services.Http.Addr = conf.Val.(string)

	conf, ok = c.keys["HTTP_BIND_PORT"]
	if ok {
		if conf.ValType != "int" {
			return fmt.Errorf("invalid value type %s for configuration key %s", conf.ValType, "HTTP_BIND_ADDR (must be an int)")
		}
	}
	c.services.Http.Port = conf.Val.(int)

	// Cookie Service
	conf, ok = c.keys["AUTH_COOKIE_NAME"]
	if ok {
		if conf.ValType != "string" {
			return fmt.Errorf("invalid value type %s for configuration key %s", conf.ValType, "AUTH_COOKIE_NAME (must be a string)")
		}
	}
	c.services.Cookie.Name = conf.Val.(string)

	conf, ok = c.keys["AUTH_COOKIE_EXPIRE"]
	if ok {
		if conf.ValType != "int" {
			return fmt.Errorf("invalid value type %s for configuration key %s", conf.ValType, "AUTH_COOKIE_EXPIRE (must be an int)")
		}
	}
	c.services.Cookie.MaxAge = conf.Val.(int)

	conf, ok = c.keys["AUTH_COOKIE_DOMAIN"]
	if ok {
		if conf.ValType != "string" {
			return fmt.Errorf("invalid value type %s for configuration key %s", conf.ValType, "AUTH_COOKIE_DOMAIN (must be a string)")
		}
	}
	c.services.Cookie.Domain = conf.Val.(string)

	conf, ok = c.keys["AUTH_COOKIE_HTTP_ONLY"]
	if ok {
		if conf.ValType != "bool" {
			return fmt.Errorf("invalid value type %s for configuration key %s", conf.ValType, "AUTH_COOKIE_HTTP_ONLY (must be a bool)")
		}
		tempBool := conf.Val.(int)
		if tempBool != 0 {
			c.services.Cookie.HttpOnly = true
		}
	}

	conf, ok = c.keys["AUTH_COOKIE_SECURE"]
	if ok {
		if conf.ValType != "bool" {
			return fmt.Errorf("invalid value type %s for configuration key %s", conf.ValType, "AUTH_COOKIE_SECURE (must be a bool)")
		}
		tempBool := conf.Val.(int)
		if tempBool != 0 {
			c.services.Cookie.Secure = true
		}
	}

	conf, ok = c.keys["AUTH_SECRET_KEY"]
	if ok {
		if conf.ValType != "string" {
			return fmt.Errorf("invalid value type %s for configuration key %s", conf.ValType, "AUTH_SECRET_KEY (must be a string)")
		}
	}
	c.services.Cookie.Domain = conf.Val.(string)

	// MemDbs
	// Permissions
	conf, ok = c.keys["DB_PERMISSIONS_ADDR"]
	if ok {
		if conf.ValType != "string" {
			return fmt.Errorf("invalid value type %s for configuration key %s", conf.ValType, "DB_PERMISSIONS_ADDR (must be a string)")
		}
	}
	c.services.MemDbs.Permissions.Addr = conf.Val.(string)

	conf, ok = c.keys["DB_PERMISSIONS_PORT"]
	if ok {
		if conf.ValType != "int" {
			return fmt.Errorf("invalid value type %s for configuration key %s", conf.ValType, "DB_PERMISSIONS_PORT (must be an int)")
		}
	}
	c.services.MemDbs.Permissions.Port = conf.Val.(int)

	conf, ok = c.keys["DB_PERMISSIONS_DB"]
	if ok {
		if conf.ValType != "int" {
			return fmt.Errorf("invalid value type %s for configuration key %s", conf.ValType, "DB_PERMISSIONS_DB (must be an int)")
		}
	}
	c.services.MemDbs.Permissions.Db = conf.Val.(int)

	// Sessions
	conf, ok = c.keys["DB_SESSIONS_ADDR"]
	if ok {
		if conf.ValType != "string" {
			return fmt.Errorf("invalid value type %s for configuration key %s", conf.ValType, "DB_SESSIONS_ADDR (must be a string)")
		}
	}
	c.services.MemDbs.Sessions.Addr = conf.Val.(string)

	conf, ok = c.keys["DB_SESSIONS_PORT"]
	if ok {
		if conf.ValType != "int" {
			return fmt.Errorf("invalid value type %s for configuration key %s", conf.ValType, "DB_SESSIONS_PORT (must be an int)")
		}
	}
	c.services.MemDbs.Sessions.Port = conf.Val.(int)

	conf, ok = c.keys["DB_SESSIONS_DB"]
	if ok {
		if conf.ValType != "int" {
			return fmt.Errorf("invalid value type %s for configuration key %s", conf.ValType, "DB_SESSIONS_DB (must be an int)")
		}
	}
	c.services.MemDbs.Sessions.Db = conf.Val.(int)

	// History
	conf, ok = c.keys["DB_HISTORY_ADDR"]
	if ok {
		if conf.ValType != "string" {
			return fmt.Errorf("invalid value type %s for configuration key %s", conf.ValType, "DB_HISTORY_ADDR (must be a string)")
		}
	}
	c.services.MemDbs.History.Addr = conf.Val.(string)

	conf, ok = c.keys["DB_HISTORY_PORT"]
	if ok {
		if conf.ValType != "int" {
			return fmt.Errorf("invalid value type %s for configuration key %s", conf.ValType, "DB_HISTORY_PORT (must be an int)")
		}
	}
	c.services.MemDbs.History.Port = conf.Val.(int)

	conf, ok = c.keys["DB_HISTORY_DB"]
	if ok {
		if conf.ValType != "int" {
			return fmt.Errorf("invalid value type %s for configuration key %s", conf.ValType, "DB_HISTORY_DB (must be an int)")
		}
	}
	c.services.MemDbs.History.Db = conf.Val.(int)

	// Logs
	conf, ok = c.keys["DB_LOGS_ADDR"]
	if ok {
		if conf.ValType != "string" {
			return fmt.Errorf("invalid value type %s for configuration key %s", conf.ValType, "DB_LOGS_ADDR (must be a string)")
		}
	}
	c.services.MemDbs.Logs.Addr = conf.Val.(string)

	conf, ok = c.keys["DB_LOGS_PORT"]
	if ok {
		if conf.ValType != "int" {
			return fmt.Errorf("invalid value type %s for configuration key %s", conf.ValType, "DB_LOGS_PORT (must be an int)")
		}
	}
	c.services.MemDbs.Logs.Port = conf.Val.(int)

	conf, ok = c.keys["DB_LOGS_DB"]
	if ok {
		if conf.ValType != "int" {
			return fmt.Errorf("invalid value type %s for configuration key %s", conf.ValType, "DB_LOGS_DB (must be an int)")
		}
	}
	c.services.MemDbs.Logs.Db = conf.Val.(int)

	return nil
}

func (c *Config) getEnvLogs() {

	tempStr := os.Getenv("SERVICE_LOG_LEVEL")
	tempInt, err := strconv.Atoi(tempStr)
	if err != nil {
		c.log.Level = defaultLogLevel

	} else if c.log.Level < LOG_LEVEL_NONE {
		c.log.Level = LOG_LEVEL_NONE
	} else if c.log.Level > LOG_LEVEL_DEBUG {
		c.log.Level = LOG_LEVEL_DEBUG
	} else {
		c.log.Level = LogLevels(tempInt)
	}

	switch c.log.Level {
	case LOG_LEVEL_DEBUG:
		slog.SetLogLoggerLevel(slog.LevelDebug)

	case LOG_LEVEL_ERROR:
		slog.SetLogLoggerLevel(slog.LevelError)

	case LOG_LEVEL_WARNING:
		slog.SetLogLoggerLevel(slog.LevelWarn)

	case LOG_LEVEL_INFO:
		slog.SetLogLoggerLevel(slog.LevelInfo)
	}
}

func (c *Config) getEnvDatabase() error {

	var err error

	c.services.Db.Host = os.Getenv("DB_HOST")
	if c.services.Db.Host == "" {
		slog.Error(ErrNoHosts.Error())
		return ErrNoHosts
	}

	c.services.Db.Port, err = strconv.Atoi(os.Getenv("DB_PORT"))
	if err != nil {
		slog.Error(ErrNoPort.Error())
		return ErrNoPort
	}

	c.services.Db.Database = os.Getenv("DB_DATABASE")
	if c.services.Db.Database == "" {
		slog.Error(ErrNoDatabaseName.Error())
		return ErrNoDatabaseName
	}

	c.services.Db.Username = os.Getenv("DB_USERNAME")
	c.services.Db.Password = os.Getenv("DB_PASSWORD")
	c.services.Db.Options = os.Getenv("DB_OPTIONS")

	return nil
}
