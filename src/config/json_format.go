package config

import (
	jsoniter "github.com/json-iterator/go"
)

type JsonConfigFormat struct {
	BaseConfigFormat
}

func NewJsonConfigFormat(fName string) *JsonConfigFormat {
	return &JsonConfigFormat{
		BaseConfigFormat: BaseConfigFormat{fileName: fName},
	}
}

func (t *JsonConfigFormat) CnfWrite(cnf *Config) error {
	return t.ReadJsonFile(cnf)
}

func (t *JsonConfigFormat) ReadJsonFile(obj interface{}) (err error) {
	buf, err := t.ReadFile(t.fileName)
	if err != nil {
		return
	}
	return jsoniter.Unmarshal(buf, obj)
}
