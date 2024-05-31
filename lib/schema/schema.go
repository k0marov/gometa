package schema

import (
	"fmt"
	"log"
	"reflect"
	"slices"
)

type Entity struct {
	Name     string
	JsonName string
	Fields   []Field
}

type Field struct {
	GoName   string
	JsonName string
	Type     FieldType
}

type FieldType string

const (
	Int      = FieldType("int")
	Uint64   = FieldType("uint64")
	Float    = FieldType("float")
	Bool     = FieldType("bool")
	String   = FieldType("string")
	TimeUnix = FieldType("time_unix")
)

const UnixTimeExample = 1694801985

func fieldTypeFromInterface(i any) (FieldType, error) {
	kind := reflect.TypeOf(i).Kind()
	switch kind {
	case reflect.Float64:
		num := i.(float64)
		if num == float64(int(num)) {
			if int(num) == UnixTimeExample {
				return TimeUnix, nil
			}
			return Int, nil
		}
		return Float, nil
	case reflect.Bool:
		return Bool, nil
	case reflect.String:
		return String, nil
	}
	return "", fmt.Errorf("unknown type %v kind, value %v i", kind, i)
}

func (e Entity) HasTimeField() bool {
	return slices.ContainsFunc(e.Fields, func(f Field) bool { return f.Type == TimeUnix })
}

func (t FieldType) GolangType() string {
	switch t {
	case Int:
		return "int"
	case Uint64:
		return "uint64"
	case Float:
		return "float64"
	case String:
		return "string"
	case TimeUnix:
		return "time.Time"
	}
	log.Panicf("Unknown field type %s", t)
	return "UNKNOWN"
}

func (f Field) IsPrimaryKey() bool {
	return f.JsonName == PrimaryKeyName
}

func (f Field) GetGormTags() string {
	if f.IsPrimaryKey() && f.Type == String {
		return "`gorm:\"default:uuid_generate_v4()\"`"
	}
	return ""
}

func (f Field) GetJsonTags() string {
	return fmt.Sprintf("`json:\"%s\"`", f.JsonName)
}
