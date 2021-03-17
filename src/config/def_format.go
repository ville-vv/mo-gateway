package config

import (
	"github.com/ville-vv/vilgo/vsql"
	"io/ioutil"
	"os"
)

type CnfWriter interface {
	CnfWrite(cnf *Config) error
}

type BaseConfigFormat struct {
	fileName string
	fType    string
}

func (b *BaseConfigFormat) ReadFile(fileName string) (data []byte, err error) {
	var (
		file *os.File
	)
	if file, err = os.Open(fileName); err != nil {
		return
	}
	defer file.Close()
	return ioutil.ReadAll(file)
}

type DefaultConfigFormat struct {
	BaseConfigFormat
}

func (t *DefaultConfigFormat) CnfWrite(cnf *Config) error {
	cnf.Mysql = &vsql.MySqlCnf{
		Version:   "8",
		UserName:  "root",
		Address:   "192.168.3.8:3306",
		Password:  "Root123.",
		Default:   "information_schema",
		MaxIdles:  10,
		MaxOpens:  1000,
		Databases: []string{"vil_user"},
	}
	cnf.Server = &ServerCnf{
		Host: "0.0.0.0",
		Port: "7001",
	}
	return nil
}
