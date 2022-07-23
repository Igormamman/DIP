package tools

import (
	"reflect"
	"strings"
)

func GetTagByField(structPoint interface{}, fieldPointer interface{}, key string) (tag string) {

	val := reflect.ValueOf(structPoint).Elem()
	val2 := reflect.ValueOf(fieldPointer).Elem()

	for i := 0; i < val.NumField(); i++ {
		valueField := val.Field(i)
		if valueField.Addr().Interface() == val2.Addr().Interface() {
			values := strings.Split(val.Type().Field(i).Tag.Get(key), ",")
			if len(values) > 0 {
				return values[0]
			} else {
				return
			}
		}
	}
	return
}
