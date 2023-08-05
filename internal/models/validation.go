package models

import (
	"github.com/go-playground/validator/v10"
	"reflect"
	"strings"
)

// Функция для кастомной валидации структуры OrderJSON
func validateNoSpecialChars(sl validator.StructLevel) {
	order := sl.Current().Interface().(OrderJSON)
	orderType := reflect.TypeOf(order)
	orderValue := reflect.ValueOf(order)

	// Определите здесь, какие символы считать специальными
	// В этом примере, запрещаем точку с запятой и двойной дефис
	for i := 0; i < orderValue.NumField(); i++ {
		fieldValue := orderValue.Field(i)
		if fieldValue.Kind() == reflect.String {
			fieldName := orderType.Field(i).Name
			fieldStr := fieldValue.String()

			for _, char := range []string{";", "--"} {
				if strings.Contains(fieldStr, char) {
					sl.ReportError(fieldValue, fieldName, "noSpecialChars", "", "")
				}
			}
		}
	}
}

// Регистрируем кастомный валидатор для структуры OrderJSON
func init() {
	validate.RegisterStructValidation(validateNoSpecialChars, OrderJSON{})
}
