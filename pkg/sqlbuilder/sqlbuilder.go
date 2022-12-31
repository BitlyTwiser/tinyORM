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
	Err              error
	Query            string
	Args             []any
	TableName        string
	model            any
	Attributes       []string
	mappedAttributes map[string]any
}

func serializeModelData(model any) *Query {
	// Enforce usage of pointer, else everything will fail
	if isPointer(model) {
		return &Query{Err: fmt.Errorf("pointer not passed. Please pass a pointer to the model")}
	}

	// Name of the struct itself, which is the DB table name
	tableName := lowerSnakeCase(reflect.TypeOf(model).Elem().Name())
	q := &Query{
		model:            model,
		TableName:        tableName,
		mappedAttributes: make(map[string]any),
	}

	// Check the pluralization of the tableName. If its not plural, pluralize it by adding s
	// ToDO: Make this less pathetic
	if !strings.HasSuffix(tableName, "s") {
		q.TableName = tableName + "s"
	}

	nVal := reflect.Indirect(reflect.ValueOf(model))

	// If slice, make higher level call deal with it.
	if nVal.Kind() == reflect.Slice {
		// We return the name of the table
		return q
	}

	// Parse attributes and values from passed in model
	for i := 0; i < nVal.NumField(); i++ {
		var name string
		f := nVal.Type().Field(i)

		// If a DB tag is present, take this field instead. Else, parse field from struct attribute
		if t, ok := f.Tag.Lookup("db"); ok {
			name = t
			q.Attributes = append(q.Attributes, t)
		} else {
			lowerSnakeName := lowerSnakeCase(f.Name)
			name = lowerSnakeName
			q.Attributes = append(q.Attributes, lowerSnakeName)
		}

		value := nVal.Field(i).Interface()

		q.mappedAttributes[name] = value
		q.Args = append(q.Args, value)
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
		queryString.WriteString(INSERT + " INTO " + q.TableName + " " + q.createTableString(databaseType))
	case DELETE:
		queryString.WriteString(DELETE + " FROM " + q.TableName + " " + q.deleteString(databaseType))
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

// Reflect the attributes from given reflect.Value and passed back slice of pointers to found attributes
// Generally to be used for destructuring a reflect.Slice type
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
func (q *Query) createTableString(databaseType string) string {
	var colString strings.Builder
	var valString strings.Builder
	var valSymbol string
	colString.WriteString("(")
	valString.WriteString("(")

	for i, v := range q.Attributes {
		// PSQL uses $ for values
		if databaseType == "psql" {
			valSymbol = "$" + strconv.Itoa(i+1)
		} else {
			valSymbol = "?"
		}
		if i != 0 {
			v = " " + v
		}
		colString.WriteString(v + ",")

		if i != 0 {
			valSymbol = " " + valSymbol
		}

		valString.WriteString(valSymbol + ",")
	}

	tmp := strings.TrimSuffix(colString.String(), ",")
	colString.Reset()
	colString.WriteString(tmp + ")")

	tmp = strings.TrimSuffix(valString.String(), ",")
	valString.Reset()
	valString.WriteString(tmp + ")")

	return (colString.String() + " VALUES " + valString.String())
}

func (q *Query) deleteString(databaseType string) string {
	var s strings.Builder
	var valSymbol string = "?" // Default to ?

	if databaseType == "psql" {
		valSymbol = "$" // Set to only first attribute as we are only going after the id
	}

	// If id is found, write query, remove all attributes for query aside from ID
	if id, found := q.mappedAttributes["id"]; found {
		q.Args = []any{id}

		s.WriteString("WHERE id = " + (valSymbol + "1"))

		return s.String()
	}

	// No ID is present, do any fields have values?
	// If not, bulk delete from table
	valSymbol = ""
	for i, attr := range q.Attributes {
		iter := strconv.Itoa(i)
		if i == 0 {
			s.WriteString(fmt.Sprintf("WHERE %s = %s", attr, (valSymbol + iter)))

			continue
		}

		s.WriteString(fmt.Sprintf("AND %s = %s", attr, (valSymbol + iter)))
	}

	return s.String()
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
