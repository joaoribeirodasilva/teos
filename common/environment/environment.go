package environment

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"strconv"

	"github.com/joaoribeirodasilva/teos/common/utils/dump"
	"github.com/joho/godotenv"
)

const (
	defaultBindIp          = "127.0.0.1"
	defaultLogLevel        = 3
	defaultDatabaseOptions = ""
	defaultEnvironmentName = "prod"
)

var (
	ErrNoDsnAndNoProtocol     = errors.New("no DB_DSN and no DB_PROTOCOL defined")
	ErrNoDsnAndNoHosts        = errors.New("no DB_DSN and no DB_HOSTS defined")
	ErrNoDsnAndNoDatabaseName = errors.New("no DB_DSN and no DB_DATABASE defined")
	ErrJSONFailed             = errors.New("failed to generate JSON string for the configuration object")

	LOG_LEVEL_NONE    = 0
	LOG_LEVEL_ERROR   = 1
	LOG_LEVEL_WARNING = 2
	LOG_LEVEL_INFO    = 3
	LOG_LEVEL_DEBUG   = 4
)

// EnvDatabase contains the values for database configuration
type EnvDatabase struct {
	Dsn      string `json:"dsn"`
	Host     string `json:"hosts"`
	Name     string `json:"name"`
	Username string `json:"username"`
	Password string `json:"-"`
	Options  string `json:"options"`
}

// EnvApplication contains the application configuration
type EnvApplication struct {
	Name            string `json:"name"`
	EnvironmentName string `json:"envName"`
	EnvironmentID   uint   `json:"envId"`
	LogLevel        int    `json:"logLevel"`
}

// Environment contains all the configuration structures
type Environment struct {
	Database    EnvDatabase    `json:"database"`
	Application EnvApplication `json:"application"`
}

// New creates a new environment object for the application
func New() *Environment {
	e := &Environment{}
	return e
}

// Read reads the environment file places in the same directory or
// if not found reads the environment variables from the system
func (e *Environment) Read() error {

	var err error

	slog.Info("reading environment variables")

	if err = godotenv.Load(); err != nil {

		slog.Warn("failed to find .env file, collecting env from real environment")

	}

	e.Application.Name = os.Getenv("SERVICE_NAME")

	tempStr := os.Getenv("SERVICE_LOG_LEVEL")
	e.Application.LogLevel, err = strconv.Atoi(tempStr)
	if err != nil {

		e.Application.LogLevel = defaultLogLevel
		slog.Warn(fmt.Sprintf("no SERVICE_LOG_LEVEL defined, defaulting to default log level: %d", defaultLogLevel))

	} else if e.Application.LogLevel < 0 {

		e.Application.LogLevel = LOG_LEVEL_NONE
		slog.Warn("SERVICE_LOG_LEVEL is defined, but smaller than 0, defaulting to 0")

	} else if e.Application.LogLevel > 4 {

		e.Application.LogLevel = LOG_LEVEL_DEBUG
		slog.Warn("SERVICE_LOG_LEVEL is defined, but greater than 4, defaulting to 4")

	}

	switch e.Application.LogLevel {

	case LOG_LEVEL_DEBUG:
		slog.SetLogLoggerLevel(slog.LevelDebug)

	case LOG_LEVEL_ERROR:
		slog.SetLogLoggerLevel(slog.LevelError)

	case LOG_LEVEL_WARNING:
		slog.SetLogLoggerLevel(slog.LevelWarn)

	case LOG_LEVEL_INFO:
		slog.SetLogLoggerLevel(slog.LevelInfo)

	}

	e.Application.EnvironmentName = os.Getenv("ENVIRONMENT")
	if e.Application.EnvironmentName == "" {
		e.Application.EnvironmentName = defaultEnvironmentName
	}

	e.Database.Dsn = os.Getenv("DB_DSN")
	if e.Database.Dsn == "" {
		e.Database.Host = os.Getenv("DB_HOST")
		if e.Database.Host == "" {

			slog.Error(ErrNoDsnAndNoHosts.Error())
			return ErrNoDsnAndNoHosts

		}

		e.Database.Name = os.Getenv("DB_DATABASE")
		if e.Database.Name == "" {

			slog.Error(ErrNoDsnAndNoDatabaseName.Error())
			return ErrNoDsnAndNoDatabaseName

		}

		e.Database.Username = os.Getenv("DB_USERNAME")
		e.Database.Password = os.Getenv("DB_PASSWORD")
		e.Database.Options = os.Getenv("DB_OPTIONS")
		if e.Database.Options == "" {

			e.Database.Options = defaultDatabaseOptions
			slog.Warn(fmt.Sprintf("no DB_OPTIONS defined, defaulting to default database options: %s", defaultDatabaseOptions))

		}
	}

	if e.Application.LogLevel >= LOG_LEVEL_DEBUG {

		j, err := dump.ToJSON(e)
		if err != nil {

			slog.Error(fmt.Sprintf("%s, Err: %s", ErrJSONFailed.Error(), err.Error()))
			return ErrJSONFailed

		}

		slog.Debug(fmt.Sprintf("environment now is: %s", j))

	}

	slog.Info("environment variables reded successfully")

	return nil
}
