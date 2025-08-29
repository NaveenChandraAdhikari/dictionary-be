package handlers

import (
	"errors"
	"fmt"
	"reflect"
	"restapi/internal/pkg/utils"
	"strings"
)

func GetFieldNames(model interface{}) []string {

	val := reflect.TypeOf(model)

	fields := []string{}

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldToAdd := strings.TrimSuffix(field.Tag.Get("json"), ",omitempty")
		fields = append(fields, fieldToAdd)

	}
	return fields
}

func CheckBlankFields(value interface{}) error {
	val := reflect.ValueOf(value)
	for i := 0; i < val.NumField(); i++ {

		field := val.Field(i)
		if field.Kind() == reflect.String && field.String() == "" {
			fmt.Println("field.Kind()", field.Kind())
			fmt.Println("reflect.String", reflect.String)
			fmt.Println("field.String()", field.String())

			//http.Error(w, "All fields are required", http.StatusBadRequest)
			return utils.ErrorHandler(errors.New("all fields are required"), "All fields are required")
		}
	}
	return nil
}
