package configbuilder

import (
	"fmt"
	"reflect"
	"sync"
)

type Config map[interface{}]interface{}

type ConfigBuilder struct {
	sync.RWMutex
	metaMap map[string]*configMeta
}

func NewConfigBuilder() *ConfigBuilder {
	return &ConfigBuilder{
		metaMap: make(map[string]*configMeta),
	}
}

func (self *ConfigBuilder) MakeConfig(fileType int, configPath string, structType interface{}) Config {
	self.RLock()
	meta := self.metaMap[configPath]
	self.RUnlock()

	if meta != nil {
		return meta.config
	}

	var parser Parser
	switch fileType {
	case FileType_Csv:
		parser = newCsvParser()
	default:
		panic("[-] ConfigBuilder: type not support")
	}

	meta = newConfigMeta(configPath)
	meta.parseStructType(parser.GetTypeString(), reflect.TypeOf(structType))

	if err := parser.GenerateConfig(meta); err != nil {
		msg := fmt.Sprintf("[-] ConfigBuilder:%s %s", configPath, err.Error())
		panic(msg)
	}

	self.Lock()
	self.metaMap[configPath] = meta
	self.Unlock()

	return meta.config
}

func (self *ConfigBuilder) Clear() {
	self.Lock()
	self.metaMap = make(map[string]*configMeta)
	self.Unlock()
}
