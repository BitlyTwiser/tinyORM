package sqlbuilder

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/google/uuid"
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
	mappedAttributes map[string]attribute
}

type attribute struct {
	value any
	t     reflect.Kind
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
		mappedAttributes: make(map[string]attribute),
	}

	// Check the pluralization of the tableName. If its not plural, pluralize it by adding s
	// ToDO: Make this less pathetic
	if !strings.HasSuffix(tableName, "s") {
		q.TableName = tableName + "s"
	}

	nVal := reflect.Indirect(reflect.ValueOf(model))

	// If slice, make higher level call deal with it.
	if nVal.Kind() == reflect.Slice {
		return q
	}

	// Parse attributes and values from passed in model
	for i := 0; i < nVal.NumField(); i++ {
		var name string
		f := nVal.Type().Field(i)

		// If a DB tag is present, take this field instead. Else, parse field from struct attribute
		// If the models value is nil or empty, the attribute is removed
		if t, ok := f.Tag.Lookup("db"); ok {
			name = t
		} else {
			name = lowerSnakeCase(f.Name)
		}

		// Not sure if this is a good idea, if we exclued fields like this, could this lead to issues?
		if value := nVal.Field(i); value.IsValid() && !value.IsZero() {
			v := value.Interface()
			q.mappedAttributes[name] = attribute{value: v, t: value.Kind()}
			q.Args = append(q.Args, v)
			q.Attributes = append(q.Attributes, name)
		}
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
		query, err := q.updateString(databaseType)

		if err != nil {
			q.Err = err

			return *q
		}
		queryString.WriteString(UPDATE + " " + q.TableName + " SET " + query)
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

	// Find and where are treated differently, just return aggregated data here.
	fwReg := regexp.MustCompile(fmt.Sprintf(`(?m)(%s|%s)`, FIND, WHERE))
	fwMatch := fwReg.Match([]byte(queryType))

	if fwMatch {
		return *serializeModelData(model)
	}

	// Nothing was found matching that string
	return Query{Err: fmt.Errorf("no matching query builder was found for the string %s", queryType)}
}

// CoalesceQueryBuilder will wrap the incoming stmt and query attributes in COALESCE with the default types per each attribute.
// This will avoid errors when null data is found when using find/where
// Uses all types and names of attributes from passed in model
func CoalesceQueryBuilder(model reflect.Type) string {
	var coalesceQuery strings.Builder
	coalesceString := " COALESCE"

	// Reflect the pointer out from the slice attributes
	//innerV := model.Type().Elem()
	for i := 0; i < model.NumField(); i++ {
		var name string
		val := model.Field(i)
		// If the models value is nil or empty, the attribute is removed
		if t, ok := val.Tag.Lookup("db"); ok {
			name = t
		} else {
			name = lowerSnakeCase(val.Name)
		}
		switch val.Type.Kind() {
		case reflect.String:
			coalesceQuery.WriteString(fmt.Sprintf(" %s(%s, %s),", coalesceString, name, "''"))
		case reflect.Array:
			// This generally would mean a jsonb array or other
			if name == "id" {
				coalesceQuery.WriteString(fmt.Sprintf(" %s(%s, %v),", coalesceString, name, "'00000000-00000000-00000000-00000000'"))

				continue
			}

			// 0 sized array
			coalesceQuery.WriteString(fmt.Sprintf(" %s(%s, '%v'),", coalesceString, name, [0]any{}))
		case reflect.Map:
			// jsonb column
			coalesceQuery.WriteString(fmt.Sprintf(" %s(%s, '{}'),", coalesceString, name))
		case reflect.Int:
			// This is the default for struct types, so should work here.
			coalesceQuery.WriteString(fmt.Sprintf(" %s(%s, %d),", coalesceString, name, 0))
		case reflect.Int8:
			coalesceQuery.WriteString(fmt.Sprintf(" %s(%s, %d),", coalesceString, name, 0))
		case reflect.Int16:
			coalesceQuery.WriteString(fmt.Sprintf(" %s(%s, %d),", coalesceString, name, 0))
		case reflect.Int32:
			coalesceQuery.WriteString(fmt.Sprintf(" %s(%s, %d),", coalesceString, name, 0))
		case reflect.Int64:
			coalesceQuery.WriteString(fmt.Sprintf(" %s(%s, %d),", coalesceString, name, 0))
		case reflect.Bool:
			coalesceQuery.WriteString(fmt.Sprintf(" %s(%s, %v),", coalesceString, name, false))
		case reflect.Float64:
			coalesceQuery.WriteString(fmt.Sprintf(" %s(%s, %f),", coalesceString, name, 0.0))
		case reflect.Float32:
			coalesceQuery.WriteString(fmt.Sprintf(" %s(%s, %f),", coalesceString, name, 0.0))
		case reflect.Interface:
			// Best guess, try string?
			coalesceQuery.WriteString(fmt.Sprintf(" %s(%s, %v),", coalesceString, name, ""))
		case reflect.Slice:
			// Any slice
			coalesceQuery.WriteString(fmt.Sprintf(" %s(%s, '%v'),", coalesceString, name, []any{}))
		}
	}

	return strings.TrimSpace(strings.TrimSuffix(coalesceQuery.String(), ","))
}

// Maps out values pulled from struct pointer and parses data into a string
// The resulting string is the query to set the values for the INSERT query
func (q *Query) createTableString(databaseType string) string {
	var colString strings.Builder
	var valString strings.Builder
	var valSymbol string
	colString.WriteString("(")
	valString.WriteString("(")

	// If ID was not passed with model record being created, generate one.
	if _, found := q.mappedAttributes["id"]; !found {
		id := uuid.New()
		q.Attributes = append(q.Attributes, "id")
		q.Args = append(q.Args, id)
		q.mappedAttributes["id"] = attribute{value: id, t: reflect.TypeOf(id).Kind()}
	}

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
	if a, found := q.mappedAttributes["id"]; found {
		q.Args = []any{a.value}

		s.WriteString("WHERE id = " + (valSymbol + "1"))

		return s.String()
	}

	// No ID is present, do any fields have values?
	// If not, bulk delete from table
	for i, attr := range q.Attributes {
		iter := strconv.Itoa(i + 1)
		if i == 0 {
			s.WriteString(fmt.Sprintf("WHERE %s = %s", attr, (valSymbol + iter)))

			continue
		}

		s.WriteString(fmt.Sprintf(" AND %s = %s", attr, (valSymbol + iter)))
	}

	return s.String()
}

func (q *Query) updateString(databaseType string) (string, error) {
	var s strings.Builder
	var valSymbol string = "?"

	if databaseType == "psql" {
		valSymbol = "$1"
	}

	if _, found := q.mappedAttributes["id"]; !found {
		return s.String(), fmt.Errorf("no id was passed, id must be present for update")
	}

	for k, v := range q.mappedAttributes {
		if k == "id" {
			continue
		}
		s.WriteString(fmt.Sprintf(" %s = '%v',", k, v.value))
	}

	tmp := strings.TrimSpace(strings.TrimSuffix(s.String(), ","))
	s.Reset()
	s.WriteString(tmp + " " + "WHERE id = " + valSymbol)

	return s.String(), nil
}

func (q *Query) GetModelID() any {
	if v, found := q.mappedAttributes["id"]; found {
		return v.value
	}

	return nil
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
