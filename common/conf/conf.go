package conf

import (
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type ConfService struct {
	Name        string
	BindAddress string
	BindPort    int
}

type ConfDatabase struct {
	Host     string
	Port     int
	Database string
	Username string
	Password string
	Options  string
}

type ConfRedis struct {
	Host     string
	Port     int
	Username string
	Password string
}

type Conf struct {
	Service         ConfService
	Database        ConfDatabase
	Redis           ConfRedis
	DefaultPageSize int64
}

const (
	defaultBindIp          = "0.0.0.0"
	defaultBindPort        = 8081
	defaultDatabasePort    = 3306
	defaultDatabaseName    = "teos"
	defaultDatabaseOptions = "charset=utf8&parseTime=True"
	defaultPageSize        = 100
	defaultRedisPort       = 6379
)

func New(serviceName string) *Conf {
	c := &Conf{}
	c.Service.Name = serviceName
	return c
}

func (c *Conf) Read() bool {

	if err := godotenv.Load(); err != nil {
		//fmt.Println(err.Error())
		slog.Info("[CONF] .env file not found, will search in the system environment variables")

	}

	tempStr := strings.TrimSpace(os.Getenv("SERVICE_NAME"))
	if tempStr == "" {
		slog.Error(fmt.Sprintf("[%s] - %s", "COMMON::CONF::Read", "invalid service name"))
		return false
	} else if c.Service.Name != tempStr {
		slog.Error(fmt.Sprintf("[%s] - %s", "COMMON::CONF::Read", "the configuration found is for service %s and this service is %s", tempStr, c.Service.Name))
		return false
	}

	c.Service.BindAddress = defaultBindIp
	tempStr = strings.TrimSpace(os.Getenv("SERVICE_IP"))
	if tempStr != "" {
		c.Service.BindAddress = tempStr
	}

	c.Service.BindPort = defaultBindPort
	tempStr = strings.TrimSpace(os.Getenv("SERVICE_PORT"))
	if tempStr != "" {
		tempInt, err := strconv.Atoi(tempStr)
		if err != nil {
			slog.Error(fmt.Sprintf("[%s] - %s", "COMMON::CONF::Read", "invalid service port found"))
			return false
		}
		c.Service.BindPort = tempInt
	}

	tempStr = strings.TrimSpace(os.Getenv("DB_HOST"))
	if tempStr == "" {
		slog.Error(fmt.Sprintf("[%s] - %s", "COMMON::CONF::Read", "invalid database host address found"))
		return false
	}
	c.Database.Host = tempStr

	c.Database.Port = defaultDatabasePort
	tempStr = strings.TrimSpace(os.Getenv("DB_PORT"))
	if tempStr != "" {
		tempInt, err := strconv.Atoi(tempStr)
		if err != nil {
			slog.Error(fmt.Sprintf("[%s] - %s", "COMMON::CONF::Read", "invalid database port found"))
			return false
		}
		c.Database.Port = tempInt
	}

	c.Database.Database = defaultDatabaseName
	tempStr = strings.TrimSpace(os.Getenv("DB_DATABASE"))
	if tempStr != "" {
		c.Database.Database = tempStr
	}

	c.Database.Username = ""
	tempStr = strings.TrimSpace(os.Getenv("DB_USERNAME"))
	if tempStr != "" {
		c.Database.Username = tempStr
	}

	c.Database.Password = ""
	tempStr = strings.TrimSpace(os.Getenv("DB_PASSWORD"))
	if tempStr != "" {
		c.Database.Password = tempStr
	}

	c.Database.Options = defaultDatabaseOptions
	tempStr = strings.TrimSpace(os.Getenv("DB_OPTIONS"))
	if tempStr != "" {
		c.Database.Options = tempStr
	}

	c.Redis.Port = defaultRedisPort
	tempStr = strings.TrimSpace(os.Getenv("REDIS_PORT"))
	if tempStr != "" {
		tempInt, err := strconv.Atoi(tempStr)
		if err != nil {
			slog.Error(fmt.Sprintf("[%s] - %s", "COMMON::CONF::Read", "invalid service port found"))
			return false
		}
		c.Redis.Port = tempInt
	}

	c.Redis.Host = strings.TrimSpace(os.Getenv("REDIS_HOST"))
	c.Redis.Username = strings.TrimSpace(os.Getenv("REDIS_USERNAME"))
	c.Redis.Password = strings.TrimSpace(os.Getenv("REDIS_PASSWORD"))

	c.DefaultPageSize = defaultPageSize

	return true
}
