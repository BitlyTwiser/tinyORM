package databuilder

import (
	"fmt"
	"log"
	"reflect"
)

// SerializeData will serialize data from any passed in model.
// The model data will be used within the insert, create, or update methods
func SerializeData(data any) error {
	t := reflect.TypeOf(data)

	if reflect.ValueOf(data).Kind() != reflect.Ptr {
		return fmt.Errorf("must pass a pointer to struct %v", t.Name())
	}

	val := reflect.ValueOf(data).Elem()
	log.Println(val)

	sFields := reflect.VisibleFields(t)

	for _, field := range sFields {
		if field.IsExported() {
			val := field.Tag.Get(field.Name)
			if val != "" {
				fmt.Printf("found value: %v", val)
			}
			fmt.Printf("field: %v is exported", field)
		}
	}

	return nil
}
