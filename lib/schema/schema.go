package schema

import (
	"log"
	"reflect"
	"time"
)

type Entity struct {
	Name   string
	Fields []Field
}

type Field struct {
	JsonName string
	Type     FieldType
}

type FieldType string

const (
	Int        = FieldType("int")
	Float      = FieldType("float")
	Bool       = FieldType("bool")
	String     = FieldType("string")
	TimeString = FieldType("time_string")
	TimeUnix   = FieldType("time_unix")
)

const UnixTimeExample = 1694801985
const StringTimeExample = time.RFC3339

func fieldTypeFromInterface(i any) FieldType {
	kind := reflect.TypeOf(i).Kind()
	switch kind {
	case reflect.Float64:
		num := i.(float64)
		if num == float64(int(num)) {
			if int(num) == UnixTimeExample {
				return TimeUnix
			}
			return Int
		}
		return Float
	case reflect.Bool:
		return Bool
	case reflect.String:
		str := i.(string)
		if str == StringTimeExample {
			return TimeString
		}
		return String
	}
	log.Fatalf("got field type %#v of unknown kind %s", i, kind.String())
	return "unknown"
}

func (t FieldType) GolangType() string {
	switch t {
	case Int:
		return "int"
	case Float:
		return "float64"
	case String:
		return "string"
	case TimeString:
		return "time.Time"
	case TimeUnix:
		return "time.Time"
	}
	log.Panicf("Unknown field type %s", t)
	return "UNKNOWN"
}