package config

import (
	"github.com/BurntSushi/toml"
)

type TomlConfigFormat struct {
	BaseConfigFormat
}

func NewTomlConfigFormat(fName string) *TomlConfigFormat {
	return &TomlConfigFormat{
		BaseConfigFormat{fileName: fName},
	}
}

func (t *TomlConfigFormat) CnfWrite(cnf *Config) error {
	return t.ReadTomlFile(cnf)
}

func (t *TomlConfigFormat) ReadTomlFile(obj interface{}) (err error) {
	buf, err := t.ReadFile(t.fileName)
	if err != nil {
		return
	}
	return toml.Unmarshal(buf, obj)
}
