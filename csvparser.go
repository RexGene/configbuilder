package configbuilder

import (
	"fmt"
	"github.com/RexGene/csvparser"
	"reflect"
	"strconv"
)

type csvParser struct {
}

func newCsvParser() *csvParser {
	return &csvParser{}
}

func (self *csvParser) GetTypeString() string {
	return "csv"
}

func (self *csvParser) GenerateConfig(meta *configMeta) error {
	configPath := meta.configPath
	data, err := csvparser.Parse(configPath)
	if err != nil {
		return err
	}

	for lineKey, row := range data {
		ptr := reflect.New(meta.configType)
		elem := ptr.Elem()
		var keyValue interface{}
		for key, index := range meta.tagRecord {
			field, ok := row[key]
			if !ok {
				panic("[-] csvParser: field name not found:" + key)
			}

			v := elem.Field(index)
			kind := v.Type().Kind()
			switch kind {
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				if string(field) == "" {
					v.SetInt(0)
				} else {
					value, err := strconv.ParseInt(string(field), 10, 64)
					if err != nil {
						msg := fmt.Sprintf("[-] csvParser: %s field <%s:%s key:%s> could not convert to int",
							configPath, key, field, lineKey)
						panic(msg)
					}
					v.SetInt(value)
				}
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				if string(field) == "" {
					v.SetUint(0)
				} else {
					value, err := strconv.ParseUint(string(field), 10, 64)
					if err != nil {
						msg := fmt.Sprintf("[-] csvParser: %s field <%s:%s key:%s> could not convert to uint",
							configPath, key, field, lineKey)
						panic(msg)
					}
					v.SetUint(value)
				}
			case reflect.String:
				v.SetString(field.Str())
			case reflect.Float32, reflect.Float64:
				if string(field) == "" {
					v.SetFloat(0)
				} else {
					value, err := strconv.ParseFloat(string(field), 64)
					if err != nil {
						msg := fmt.Sprintf("[-] csvParser: %s field <%s:%s key:%s> could not convert to float",
							configPath, key, field, lineKey)
						panic(msg)
					}
					v.SetFloat(value)
				}
			case reflect.Bool:
				if string(field) == "" {
					v.SetBool(false)
				} else {
					value, err := strconv.ParseInt(string(field), 10, 64)
					if err != nil {
						msg := fmt.Sprintf("[-] csvParser: %s field <%s:%s key:%s> could not convert to bool",
							configPath, key, field, lineKey)
						panic(msg)
					}
					v.SetBool(value != 0)
				}
			}

			if index == meta.keyIndex {
				if !v.CanInterface() {
					panic("[-] csvParser: value could not be interface")
				}
				keyValue = v.Interface()
			}

		}

		meta.config[keyValue] = ptr.Interface()
	}

	return nil
}
