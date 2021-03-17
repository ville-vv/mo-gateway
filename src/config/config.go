package config

import (
	"fmt"
	"github.com/ville-vv/vilgo/vsql"
	"sync"
)

var (
	conf *Config
)

func Init(args ...string) *Config {
	(&sync.Once{}).Do(func() {
		switch len(args) {
		case 1:
			conf = NewConfig(args[0])
		default:
			conf = NewConfig("./conf/config.toml")
		}
	})
	return conf
}

type ServerCnf struct {
	Host string `json:"host"`
	Port string `json:"port"`
}

// Config 配置信息，可以是 toml 文件也可以是 json 文件
type Config struct {
	Server     *ServerCnf     `json:"server"`
	GrpcServer *ServerCnf     `json:"grpc_server"`
	Mysql      *vsql.MySqlCnf `json:"mysql_cnf"`
}

func NewConfig(fileName string) *Config {
	factory := NewFactoryConfig()
	return factory.GetConfig(fileName)
}

func GetMySqlCnf() *vsql.MySqlCnf {
	return conf.Mysql
}

func GetSqlDefultDb() string {
	return conf.Mysql.Default
}

func GetDbNameList() []string {
	return conf.Mysql.Databases
}

func GetDbName() string {
	return conf.Mysql.Databases[0]
}

func ServerAddress() string {
	return fmt.Sprintf("%s:%s", conf.Server.Host, conf.Server.Port)
}
