package xconfig

import (
	"os"
	"reflect"
	"strconv"

	"github.com/olvrng/rbot/be/pkg/l"
)

var ll = l.New()
var ls = ll.Sugar()

type EnvMap map[string]interface{}

func (m EnvMap) MustLoad() {
	for k, v := range m {
		MustLoadEnv(k, v)
	}
}

func MustLoadEnv(env string, val interface{}) {
	s := os.Getenv(env)
	if s == "" {
		return
	}

	v := reflect.ValueOf(val)
	if v.Kind() != reflect.Ptr {
		ls.Panicf("expect pointer for env: %v", env)
	}

	v = v.Elem()
	switch v.Kind() {
	case reflect.String:
		v.SetString(s)
	case reflect.Bool:
		b, err := strconv.ParseBool(s)
		if err != nil {
			ls.Fatalf("%v expects a boolean, got: %v", env, s)
		}
		v.SetBool(b)
	case reflect.Int, reflect.Int64:
		i, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			ls.Fatalf("%v expects an integer, got: %v", env, s)
		}
		v.SetInt(i)
	case reflect.Uint, reflect.Uint64:
		i, err := strconv.ParseUint(s, 10, 64)
		if err != nil {
			ls.Fatalf("%v expects an unsigned integer, got: %v", env, s)
		}
		v.SetUint(i)
	case reflect.Float64:
		i, err := strconv.ParseFloat(s, 64)
		if err != nil {
			ls.Fatalf("%v expects an float64, got: %v", env, s)
		}
		v.SetFloat(i)
	default:
		ls.Panicf("unexpected type for env: %v, type: %v", env, v.Kind())
	}
}
