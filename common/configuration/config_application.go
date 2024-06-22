package configuration

type EnvironmentTypes int

const (
	defaultEnvironmentID = 5

	EnvDev EnvironmentTypes = iota + 1
	EnvLocalDev
	EnvTest
	EnvEval
	EnvProd
)

var (
	EnvironmentKeys = map[string]EnvironmentTypes{
		"dev":   EnvDev,
		"local": EnvLocalDev,
		"test":  EnvTest,
		"eval":  EnvEval,
		"prod":  EnvProd,
	}

	EnvironmentNames = map[EnvironmentTypes]string{
		EnvDev:      "dev",
		EnvLocalDev: "local",
		EnvTest:     "test",
		EnvEval:     "eval",
		EnvProd:     "prod",
	}
)

type ConfigApplication struct {
	Code           string
	ID             uint
	EnvironmentKey string
	EnvironmentID  EnvironmentTypes
}
