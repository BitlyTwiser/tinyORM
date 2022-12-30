package sqlbuilder

import (
	"fmt"
	"reflect"
	"strings"
)

// Standard data manupulation operations
const (
	SELECT = "select"
	CREATE = "create"
	INSERT = "insert"
	UPDATE = "update"
	DELETE = "delete"
	WHERE  = "select"
)

type Query struct {
	Err   error
	Query string
	Args  []any
}

// SerializeData will serialize data from any passed in model.
// The model data will be used within the insert, create, or update methods
// Does NOT fill the model with any data
func serializeData(model any) error {
	t := reflect.TypeOf(model)

	if reflect.ValueOf(model).Kind() != reflect.Ptr {
		return fmt.Errorf("must pass a pointer to struct %v", t.Name())
	}

	val := reflect.TypeOf(model).Elem()
	val.Name()

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

// Will query for data and adjust the model accordingly with the foudn data values.
// Errors will occur when any inproper models are passed.
// In theory, if the qwuery successed, this process should work without issue.
func serializeAndModify(model any, data any) error {
	// Look at using Elem().Set() to set X to Y value
	// May work to set values after parsing query row data
	return nil
}

// Generates SQL query from given model.
// All attributes of the model are passed back as arguments to the calling function
// Name of model is lowercased, then snake cased to adhere to SQL naming conventions.
// A table is expected to exist with the given model name.
// Used for Create, Update, and Delete
func QueryBuilder(queryType string, model any) Query {
	var queryString strings.Builder
	var columnArr []string
	q := Query{}

	// Enforce usage of pointer, else everything will fail
	if reflect.ValueOf(model).Kind() != reflect.Ptr {
		return Query{Err: fmt.Errorf("pointer not passed. Please pass a pointer to the model")}
	}

	// Name of the struct itself, which is the DB table name
	tableName := lowerSnakeCase(reflect.TypeOf(model).Elem().Name())

	// Check the pluralization of the tableName. If its not plural, pluralize it by adding s

	nVal := reflect.ValueOf(model).Elem()
	// Parse attributes and values from passed in model
	for i := 0; i < nVal.NumField(); i++ {
		columnArr = append(columnArr, lowerSnakeCase(nVal.Type().Field(i).Name))
		q.Args = append(q.Args, nVal.Field(i).Interface())
	}

	switch queryType {
	case CREATE:
		queryString.WriteString(INSERT + " INTO " + tableName + " " + createTableString(columnArr))
	case DELETE:
		queryString.WriteString(DELETE + " FROM " + tableName)
	case UPDATE:
		// This could be a problem
		queryString.WriteString(UPDATE + tableName + " SET ")
	}

	q.Query = queryString.String()

	return q
}

// Maps out builder values pulled from struct pointer and parses data into a string
func createTableString(columnValues []string) string {
	var colString strings.Builder
	var valString strings.Builder
	colString.WriteString("(")
	valString.WriteString("(")

	for _, v := range columnValues {
		colString.WriteString(v + ",")
		valString.WriteString("?" + ",")
	}

	// Icing on cake
	tmp := strings.TrimSuffix(colString.String(), ",")
	colString.Reset()
	colString.WriteString(tmp + ")")

	tmp = strings.TrimSuffix(valString.String(), ",")
	valString.Reset()
	valString.WriteString(tmp + ")")

	return (colString.String() + " VALUES " + valString.String())
}

// Lowercases and Snakecases the given string as to be used in the SQL query
// Keep in mind, non Acii chars (as they can be multiple bytes in length) will not work
// Additionally, this is not true snake casing. The first char is lowered
func lowerSnakeCase(val string) string {
	var s strings.Builder

	for i, char := range val {
		// Falls in uppercase range
		if i == 0 {
			// If caps, lower case
			if isCap(char) {
				char = char + 32
			}

			s.WriteRune(char) // lower case first, but do not append underscore
			continue
		}
		if isCap(char) {
			// All runes are 32 bits apart
			// 32 bit forward run char is lowercase.
			// Saves on perform string conversion than ToLower
			s.WriteRune(char + 32)
			s.WriteString("_")
		} else {
			s.WriteRune(char)
		}
	}

	// Just in case a the last char was snake case
	return strings.TrimSuffix(s.String(), "_")
}

func isCap(char rune) bool {
	return (char >= 65 && char <= 90)
}

// Used for Find & Update queries
// Will build, query, parse, and load data aggregated from sql call into the respective model
func QueryAndUpdate(queryType string, model any, stmt string, limit int, args ...any) error {
	return nil
}
