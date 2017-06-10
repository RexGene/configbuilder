package configbuilder

import (
	"reflect"
	"strings"
)

const (
	ATTR_NAME = "attr"
	KEY_FLAG  = "key"
)

type flagValue map[string]bool

type flags map[string]flagValue

type configMeta struct {
	tagRecord  map[string]int
	configType reflect.Type
	config     Config
	configPath string
	keyIndex   int
	fieldSet   map[int]flags
}

func newConfigMeta(configPath string) *configMeta {
	return &configMeta{
		tagRecord:  make(map[string]int),
		fieldSet:   make(map[int]flags),
		config:     make(Config),
		configPath: configPath,
	}
}

func (meta *configMeta) parseStructType(fileType string, t reflect.Type) {
	kind := t.Kind()
	if kind != reflect.Ptr && kind != reflect.Struct {
		panic("[-] configMeta: type must be ptr or struct")
	}

	if kind == reflect.Ptr {
		t = t.Elem()
		if t.Kind() != reflect.Struct {
			panic("[-] configMeta: value not a struct")
		}
	}

	meta.configType = t
	tagRecord := meta.tagRecord
	for i := 0; i < t.NumField(); i++ {
		sf := t.Field(i)
		name, ok := sf.Tag.Lookup(fileType)
		if !ok {
			name = sf.Name
		} else {
			name = meta._parseOptionFields(i, name)
		}

		tagRecord[name] = i
	}
}

func (meta *configMeta) _paresFlagValue(fvStr string) flagValue {
	rs := make(flagValue)
	values := strings.Split(fvStr, ",")
	for _, value := range values {
		rs[strings.TrimSpace(value)] = true
	}

	return rs
}

func (meta *configMeta) _parseOptionFields(index int, fieldStr string) (result string) {
	fields := strings.Split(fieldStr, ";")
	f := make(flags)
	meta.fieldSet[index] = f
	for _, field := range fields {
		field = strings.TrimSpace(field)

		desc := strings.Split(field, "=")
		descLen := len(desc)
		switch descLen {
		case 1:
			meta.tagRecord[field] = index
			result = field
		case 2:
			name := strings.TrimSpace(desc[0])
			f[name] = meta._paresFlagValue(strings.TrimSpace(desc[1]))
		default:
			panic("[-] configmeta: tags format error")
		}
	}

	attrSet := f[ATTR_NAME]
	if attrSet != nil && attrSet[KEY_FLAG] {
		meta.keyIndex = index
	}

	return
}
