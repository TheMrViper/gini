package gini

import (
	"reflect"
	"strconv"
)

func ReadConfig(name string, i interface{}) error {

	resultLines, err := readLines(name)
	if err != nil {
		return err
	}
	resultMap := fileToMap(resultLines)

	mapToStruct("", resultMap, reflect.ValueOf(i).Elem())

	return nil
}

func WriteConfig(name string, i interface{}) error {
	resultMap := make(map[string]map[string]string)

	structToMap("", resultMap, reflect.ValueOf(i).Elem())

	resultLines := mapToFile(resultMap)

	return writeLines(resultLines, name)
}
func mapToStruct(name string, data map[string]map[string]string, value reflect.Value) {
	if name != "" && data[name] == nil {
		return
	}

	for f := 0; f < value.NumField(); f++ {
		valueField := value.Field(f)
		typeField := value.Type().Field(f)

		fieldName := typeField.Tag.Get("ini")
		if fieldName == "-" {
			continue
		}
		if fieldName == "" {
			fieldName = typeField.Name
		}

		fieldValue, ok := data[name][fieldName]
		if !ok {
			continue
		}

		switch valueField.Kind() {
		case reflect.Struct:
			mapToStruct(fieldName, data, valueField)
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			number, err := strconv.ParseInt(fieldValue, 10, 64)
			if err != nil {
				continue
			}

			if !valueField.OverflowInt(number) {
				valueField.SetInt(number)
			}
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			number, err := strconv.ParseUint(fieldValue, 10, 64)
			if err != nil {
				continue
			}

			if !valueField.OverflowUint(number) {
				valueField.SetUint(number)
			}
		case reflect.Float32, reflect.Float64:
			number, err := strconv.ParseFloat(fieldValue, 64)
			if err != nil {
				continue
			}

			if !valueField.OverflowFloat(number) {
				valueField.SetFloat(number)
			}
		case reflect.String:
			valueField.SetString(fieldValue)
		case reflect.Bool:
			boolean, err := strconv.ParseBool(fieldValue)
			if err != nil {
				continue
			}
			valueField.SetBool(boolean)
		}

	}
}

func structToMap(name string, data map[string]map[string]string, value reflect.Value) {
	if name != "" && data[name] == nil {
		data[name] = make(map[string]string)
	}

	for f := 0; f < value.NumField(); f++ {
		valueField := value.Field(f)
		typeField := value.Type().Field(f)

		fieldName := typeField.Tag.Get("ini")
		if fieldName == "-" {
			continue
		}
		if fieldName == "" {
			fieldName = typeField.Name
		}

		switch valueField.Kind() {
		case reflect.Struct:
			structToMap(fieldName, data, valueField)
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			data[name][fieldName] = strconv.FormatInt(valueField.Int(), 10)
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			data[name][fieldName] = strconv.FormatUint(valueField.Uint(), 10)
		case reflect.Float32, reflect.Float64:
			data[name][fieldName] = strconv.FormatFloat(valueField.Float(), 'f', 6, 64)
		case reflect.String:
			data[name][fieldName] = valueField.String()
		case reflect.Bool:
			data[name][fieldName] = strconv.FormatBool(valueField.Bool())
		}
	}
}
