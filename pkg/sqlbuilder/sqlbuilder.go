package sqlbuilder

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

// Standard data manupulation operations
const (
	SELECT = "select"
	CREATE = "create"
	INSERT = "insert"
	UPDATE = "update"
	DELETE = "delete"
	WHERE  = "where"
	FIND   = "find"
)

type Query struct {
	Err        error
	Query      string
	Args       []any
	TableName  string
	model      any
	Attributes []string
}

func serializeModelData(model any) *Query {
	// Enforce usage of pointer, else everything will fail
	if isPointer(model) {
		return &Query{Err: fmt.Errorf("pointer not passed. Please pass a pointer to the model")}
	}

	// Name of the struct itself, which is the DB table name
	tableName := lowerSnakeCase(reflect.TypeOf(model).Elem().Name())
	q := &Query{model: model, TableName: tableName}

	// Check the pluralization of the tableName. If its not plural, pluralize it by adding s
	// ToDO: Make this less pathetic
	if !strings.HasSuffix(tableName, "s") {
		q.TableName = tableName + "s"
	}

	nVal := reflect.Indirect(reflect.ValueOf(model))

	if nVal.Kind() == reflect.Slice {
		// If slice, make higher level call deal with it.
		// We return the name of the table
		return q
		// // Get inner type value
		// fmt.Println(nVal.Type().Elem())
		// fmt.Println(nVal.Type().Elem().NumField())
		// nVal = reflect.Indirect(reflect.ValueOf(nVal.Type().Elem()))
	}

	// Parse attributes and values from passed in model
	for i := 0; i < nVal.NumField(); i++ {
		f := nVal.Type().Field(i)

		// If a DB tag is present, take this field instead. Else, parse field from struct attribute
		if t, ok := f.Tag.Lookup("db"); ok {
			q.Attributes = append(q.Attributes, t)
		} else {
			q.Attributes = append(q.Attributes, lowerSnakeCase(f.Name))
		}
		q.Args = append(q.Args, nVal.Field(i).Interface())
	}

	return q
}

func (q *Query) buildQueryFromModelData(queryType string, databaseType string) Query {
	var queryString strings.Builder

	if q.Err != nil {
		return *q
	}

	switch queryType {
	case CREATE:
		queryString.WriteString(INSERT + " INTO " + q.TableName + " " + createTableString(q.Attributes, databaseType))
	case DELETE:
		queryString.WriteString(DELETE + " FROM " + q.TableName)
	case UPDATE:
		queryString.WriteString(UPDATE + q.TableName + " SET ")
	}

	q.Query = queryString.String()

	return *q
}

func (q *Query) ModelAttributes() []any {
	var pointers []any

	vals := reflect.ValueOf(q.model).Elem()
	for i := 0; i < vals.NumField(); i++ {
		pointers = append(pointers, vals.Field(i).Addr().Interface())
	}

	return pointers
}

func PointerAttributes(model reflect.Value) []any {
	var pointers []any

	model = reflect.Indirect(model)
	for i := 0; i < model.NumField(); i++ {
		pointers = append(pointers, model.Field(i).Addr().Interface())
	}

	return pointers
}

func isPointer(model any) bool {
	return reflect.ValueOf(model).Kind() != reflect.Ptr
}

// Generates SQL query from given model.
// All attributes of the model are passed back as arguments to the calling function
// Name of model is lowercased, then snake cased to adhere to SQL naming conventions.
// A table is expected to exist with the given model name.
// Used for Create, Update, and Delete
func QueryBuilder(queryType string, model any, databaseType string) Query {
	// Regex the query type to determine which pathway the function call goes
	re := regexp.MustCompile(fmt.Sprintf(`(?m)(%s|%s|%s)`, CREATE, UPDATE, DELETE))
	match := re.Match([]byte(queryType))

	if match {
		return serializeModelData(model).buildQueryFromModelData(queryType, databaseType)
	}

	fwReg := regexp.MustCompile(fmt.Sprintf(`(?m)(%s|%s)`, FIND, WHERE))
	fwMatch := fwReg.Match([]byte(queryType))

	if fwMatch {
		return *serializeModelData(model)
	}

	// Nothing was found matching that string
	return Query{Err: fmt.Errorf("no matching query builder was found for the string %s", queryType)}
}

// Used for Find & Where queries
// Will build, query, parse, and load data aggregated from sql call into the respective model
func QueryAndUpdate(queryType string, model any, stmt string, limit int, args ...any) error {
	return nil
}

// Maps out values pulled from struct pointer and parses data into a string
// The resulting string is the query to set the values for the INSERT query
func createTableString(columnValues []string, databaseType string) string {
	var colString strings.Builder
	var valString strings.Builder
	var questionVal string
	colString.WriteString("(")
	valString.WriteString("(")

	for i, v := range columnValues {
		if databaseType == "psql" {
			questionVal = "$" + strconv.Itoa(i+1)
		} else {
			questionVal = "?"
		}
		if i != 0 {
			v = " " + v
		}
		colString.WriteString(v + ",")

		if i != 0 {
			questionVal = " " + questionVal
		}

		valString.WriteString(questionVal + ",")
	}

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
			// Saves time on having to perform string conversion than ToLower function call.
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
