package config

type EnvConfigFormat struct {
	BaseConfigFormat
}

func NewEnvConfigFormat(fName string) *EnvConfigFormat {
	return &EnvConfigFormat{
		BaseConfigFormat: BaseConfigFormat{fileName: fName},
	}
}

func (t *EnvConfigFormat) CnfWrite(cnf *Config) error {
	return nil
}
