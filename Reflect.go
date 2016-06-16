package gini

import (
	"errors"
	"reflect"
	"strconv"
)

func WriteConfig(name string, i interface{}) error {
	resultMap := make(map[string]map[string]string)

	rv := reflect.ValueOf(i)
	if rv.Kind() != reflect.Ptr {
		return errors.New("Interface is not ptr")
	}
	structToMap(rv.Elem().Type().Name(), resultMap, rv.Elem())

	resultLines := mapToFile(resultMap)

	return writeLines(resultLines, name)
}
func ReadConfig(name string, i interface{}) error {

	resultLines, err := readLines(name)
	if err != nil {
		return err
	}
	resultMap := fileToMap(resultLines)

	rv := reflect.ValueOf(i)
	if rv.Kind() != reflect.Ptr {
		return errors.New("Interface is not ptr")
	}

	mapToStruct(rv.Elem().Type().Name(), resultMap, rv.Elem())

	return nil
}

func mapToStruct(name string, data map[string]map[string]string, structElem reflect.Value) {

	typeStruct := structElem.Type()

	for f := 0; f < typeStruct.NumField(); f++ {
		valueField := structElem.Field(f)
		typeField := typeStruct.Field(f)

		fieldName := typeField.Tag.Get("ini")
		if fieldName == "-" {
			continue
		}
		if fieldName == "" {
			fieldName = typeField.Name
		}

		if valueField.Kind() == reflect.Struct {
			mapToStruct(fieldName, data, valueField)
			continue
		}

		section, ok := data[name]
		if !ok {
			continue
		}
		fieldValue, ok := section[fieldName]
		if !ok {
			continue
		}

		switch valueField.Kind() {
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

func structToMap(name string, data map[string]map[string]string, structElem reflect.Value) {
	if name != "" && data[name] == nil {
		data[name] = make(map[string]string)
	}
	typeStruct := structElem.Type()

	for f := 0; f < typeStruct.NumField(); f++ {
		valueField := structElem.Field(f)
		typeField := typeStruct.Field(f)

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
