package configuration

type ConfigMemDb struct {
	Addr     string
	Port     int
	Db       int
	Username string
	Password string
}

type ConfigMemDbs struct {
	Permissions ConfigMemDb
	Sessions    ConfigMemDb
	History     ConfigMemDb
	Logs        ConfigMemDb
}
