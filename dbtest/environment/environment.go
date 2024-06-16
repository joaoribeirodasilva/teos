package environment

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"slices"
	"strconv"

	"github.com/joaoribeirodasilva/teos/dbtest/utils/dump"
	"github.com/joho/godotenv"
)

const (
	defaultBindIp          = "127.0.0.1"
	defaultBindPort        = 8080
	defaultDatabaseOptions = "retryWrites=true&w=majority"
	defaultLogLevel        = 3
)

var (
	ErrNoDsnAndNoProtocol     = errors.New("no DB_DSN and no DB_PROTOCOL defined")
	ErrNoDsnAndNoHosts        = errors.New("no DB_DSN and no DB_HOSTS defined")
	ErrNoDsnAndNoDatabaseName = errors.New("no DB_DSN and no DB_DATABASE defined")
	ErrJSONFailed             = errors.New("failed to generate JSON string for the configuration object")

	validProtocols = []string{"mongodb+srv", "mongodb"}

	LOG_LEVEL_NONE    = 0
	LOG_LEVEL_ERROR   = 1
	LOG_LEVEL_WARNING = 2
	LOG_LEVEL_INFO    = 3
	LOG_LEVEL_DEBUG   = 4
)

// EnvDatabase conatains the values for database configuration
type EnvDatabase struct {
	Dsn      string `json:"dsn"`
	Protocol string `json:"protocol"`
	Hosts    string `json:"hosts"`
	Name     string `json:"name"`
	Username string `json:"username"`
	Password string `json:"-"`
	Options  string `json:"options"`
}

// EnvApplication conatains the application configuration
type EnvApplication struct {
	Name       string `json:"name"`
	ListenIp   string `json:"bindIp"`
	ListenPort int    `json:"bindPort"`
	LogLevel   int    `json:"logLevel"`
}

// Environment conatains all the configuration structures
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

	slog.Info("reading evironment variables")

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
		slog.Warn("SERVICE_LOG_LEVEL is defined, but smalled than 0, defaulting to 0")
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

	e.Application.ListenIp = os.Getenv("SERVICE_BIND_IP")
	if e.Application.ListenIp == "" {
		e.Application.ListenIp = defaultBindIp
		slog.Warn(fmt.Sprintf("no SERVICE_BIND_IP defined, defaulting to default bind ip: %s", defaultBindIp))
	}

	tempStr = os.Getenv("SERVICE_BIND_PORT")
	e.Application.ListenPort, err = strconv.Atoi(tempStr)
	if err != nil {
		slog.Warn(fmt.Sprintf("no SERVICE_BIND_PORT defined, defaulting to default bind port: %d", defaultBindPort))
	}
	if e.Application.ListenPort < 8000 {
		e.Application.ListenPort = defaultBindPort
		slog.Warn(fmt.Sprintf("no SERVICE_BIND_PORT defined, defaulting to default bind port: %d", defaultBindPort))
	}

	e.Database.Dsn = os.Getenv("DB_DSN")
	if e.Database.Dsn == "" {
		e.Database.Protocol = os.Getenv("DB_PROTOCOL")
		if e.Database.Protocol == "" || !slices.Contains(validProtocols, e.Database.Protocol) {
			slog.Error(ErrNoDsnAndNoProtocol.Error())
			return ErrNoDsnAndNoProtocol
		}
		e.Database.Hosts = os.Getenv("DB_HOSTS")
		if e.Database.Hosts == "" {
			slog.Error(ErrNoDsnAndNoHosts.Error())
			return ErrNoDsnAndNoHosts
		}
		e.Database.Name = os.Getenv("DB_DATABASE")
		if e.Database.Name == "" {
			slog.Error(ErrNoDsnAndNoDatabaseName.Error())
			return ErrNoDsnAndNoDatabaseName
		}
		e.Database.Username = os.Getenv("DB_USERNAME")
		e.Database.Username = os.Getenv("DB_PASSWORD")
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

	slog.Info("evironment variables readed successfully")

	return nil
}
