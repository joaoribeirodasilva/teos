package configuration

type LogLevels int

const (
	LOG_LEVEL_NONE LogLevels = iota
	LOG_LEVEL_ERROR
	LOG_LEVEL_WARNING
	LOG_LEVEL_INFO
	LOG_LEVEL_DEBUG
)

type ConfigLog struct {
	Level LogLevels
}
