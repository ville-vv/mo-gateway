package config

import (
	"fmt"
	"testing"
)

func TestFactoryConfig_GetType(t *testing.T) {
	factory := &FactoryConfig{}
	fmt.Println(factory.GetType("xxx.json"))
}
