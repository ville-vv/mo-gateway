package config

import (
	"strings"
	"vilgo/vutil"
)

const (
	OriginTypeConfigEnv  = "env"
	OriginTypeConfigToml = "toml"
	OriginTypeConfigJson = "json"
	OriginTypeConfigYaml = "yaml"
	OriginTypeConfigDef  = "default"
)

type FactoryConfig struct {
	Write CnfWriter
	cnf   *Config
}

func NewFactoryConfig(wc ...CnfWriter) *FactoryConfig {
	f := new(FactoryConfig)
	if len(wc) > 0 {
		f.Write = wc[0]
	}
	return f
}

func (f *FactoryConfig) GetConfig(origin string) *Config {
	if f.Write == nil {
		format := f.GetType(origin)
		switch format {
		case OriginTypeConfigJson:
			f.Write = NewJsonConfigFormat(origin)
		case OriginTypeConfigYaml:
		case OriginTypeConfigToml:
			f.Write = NewTomlConfigFormat(origin)
		case OriginTypeConfigEnv:
			f.Write = NewEnvConfigFormat(origin)
		default:
			f.Write = &DefaultConfigFormat{}
		}
	}
	return f.newConfig()
}

func (f *FactoryConfig) newConfig() *Config {
	f.cnf = new(Config)
	if err := f.Write.CnfWrite(f.cnf); err != nil {
		panic(err)
	}
	return f.cnf
}

// GetType 配置文件类型判断
func (f *FactoryConfig) GetType(origin string) string {
	// 优先判断非文件
	switch strings.ToLower(origin) {
	case "":
		return OriginTypeConfigDef
	case OriginTypeConfigEnv:
		return OriginTypeConfigEnv
	}
	arr := strings.Split(origin, ".")
	if len(arr) > 1 {
		return arr[len(arr)-1]
	}

	// 判断是否文件存在,存在默认使用 toml
	if in := vutil.PathExists(origin); in {
		return OriginTypeConfigToml
	}
	return OriginTypeConfigDef
}
