package configbuilder

import (
	"github.com/RexGene/csvparser"
	"reflect"
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
	data, err := csvparser.Parse(meta.configPath)
	if err != nil {
		return err
	}

	for _, row := range data {
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
				v.SetInt(int64(field.Int(0)))
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				v.SetUint(uint64(field.Uint(0)))
			case reflect.String:
				v.SetString(field.Str())
			case reflect.Float32, reflect.Float64:
				v.SetFloat(float64(field.Float(0)))
			case reflect.Bool:
				v.SetBool(field.Int(0) != 0)
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
