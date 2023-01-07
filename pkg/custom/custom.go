package custom

import (
	"database/sql/driver"
	"encoding/json"
)

type Map map[string]any

// Adds Value to map
func (v *Map) Add(key string, value any) {
	(*v)[key] = value
}

// Deletes value if desired
func (v *Map) Delete(key string) {
	delete((*v), key)
}

// Returns Values to caller
func (v *Map) Values() Map {
	return *v
}

// Scanner/Valuer interface implementation
func (v *Map) Scan(value interface{}) error {
	switch value := value.(type) {
	case []byte:
		return json.Unmarshal(value, &v)
	case string:
		return json.Unmarshal([]byte(value), &v)
	default:
		return json.Unmarshal(value.([]byte), &v)
	}
}

func (v Map) Value() (driver.Value, error) {
	return json.Marshal(v)
}

func NewMap() Map {
	return make(Map)
}

type Slice []any

func (s *Slice) Append(values ...any) {
	*s = append(*s, values...)
}

// Returns Values to caller
func (s *Slice) Values() Slice {
	return *s
}

// Scanner/Valuer interface implementation
func (s *Slice) Scan(value interface{}) error {
	switch value := value.(type) {
	case []byte:
		return json.Unmarshal(value, &s)
	case string:
		return json.Unmarshal([]byte(value), &s)
	default:
		return json.Unmarshal(value.([]byte), &s)
	}
}

func (s Slice) Value() (driver.Value, error) {
	return json.Marshal(s)
}
