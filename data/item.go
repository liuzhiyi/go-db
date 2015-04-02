package data

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"time"
)

var errNilPtr = errors.New("destination pointer is nil")

type Item struct {
	raw  []interface{}
	data map[string]interface{}
}

func (i *Item) Init() {
	i.data = make(map[string]interface{})
}

func (i *Item) SetData(key string, value interface{}) {
	i.data[key] = value
}

func (i *Item) GetData(key string) interface{} {
	if val, has := i.data[key]; has {
		return val
	}
	return nil
}

func (i *Item) GetMap() map[string]interface{} {
	return i.data
}

func (i *Item) UnsetData(keys ...string) {
	if len(keys) < 1 {
		i.data = make(map[string]interface{})
	} else {
		for j := 0; j < len(keys); j++ {
			delete(i.data, keys[j])
		}
	}
}

func (i *Item) ToJson() string {
	str, _ := json.Marshal(i.data)
	return string(str)
}

func (i *Item) SetRaw(raw []interface{}) {
	i.raw = raw
}

func (i *Item) ToArray() []string {
	var row []string

	for _, val := range i.raw {
		var dst string
		i.convert(&dst, val)
		row = append(row, dst)
	}

	return row
}

func (i *Item) GetInt(key string) int {
	return int(i.GetInt64(key))
}

func (i *Item) GetInt64(key string) int64 {
	var val int64
	if i.GetData(key) == nil {
		return 0
	}
	err := i.convert(&val, i.GetData(key))
	if err != nil {
		panic(err.Error())
	}
	return val
}

func (i *Item) GetFloat64(key string) float64 {
	var val float64
	if i.GetData(key) == nil {
		return 0.0
	}
	err := i.convert(&val, i.GetData(key))
	if err != nil {
		panic(err.Error())
	}
	return val
}

func (i *Item) GetDate(key string, format string) string {
	var val string
	timeval := i.GetInt64(key)
	t := time.Unix(timeval, 0)
	val = t.Format(format)
	return val
}

func (i *Item) GetBool(key string) bool {
	var val bool
	if i.GetData(key) == nil {
		return false
	}
	err := i.convert(&val, i.GetData(key))
	if err != nil {
		panic(err.Error())
	}
	return val
}

func (i *Item) GetString(key string) string {
	var val string
	if i.GetData(key) == nil {
		return ""
	}
	err := i.convert(&val, i.GetData(key))
	if err != nil {
		panic(err.Error())
	}
	return val
}

func (i *Item) convert(dest, src interface{}) error {
	if v, ok := src.(*interface{}); ok {
		src = *v
	}
	switch s := src.(type) {
	case string:
		switch d := dest.(type) {
		case *string:
			if d == nil {
				return errNilPtr
			}
			*d = s
			return nil
		case *[]byte:
			if d == nil {
				return errNilPtr
			}
			*d = []byte(s)
			return nil
		}
	case []byte:
		switch d := dest.(type) {
		case *string:
			if d == nil {
				return errNilPtr
			}
			*d = string(s)
			return nil
		case *interface{}:
			if d == nil {
				return errNilPtr
			}
			*d = i.cloneBytes(s)
			return nil
		case *[]byte:
			if d == nil {
				return errNilPtr
			}
			*d = i.cloneBytes(s)
			return nil
		}
	case nil:
		switch d := dest.(type) {
		case *interface{}:
			if d == nil {
				return errNilPtr
			}
			*d = nil
			return nil
		case *[]byte:
			if d == nil {
				return errNilPtr
			}
			*d = nil
			return nil
		}
	}
	var sv reflect.Value
	switch d := dest.(type) {
	case *string:
		sv = reflect.ValueOf(src)
		switch sv.Kind() {
		case reflect.Bool,
			reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
			reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
			reflect.Float32, reflect.Float64:
			*d = i.asString(src)
			return nil
		}
	case *[]byte:
		sv = reflect.ValueOf(src)
		if b, ok := i.asBytes(nil, sv); ok {
			*d = b
			return nil
		}
	case *bool:
		bv, err := driver.Bool.ConvertValue(src)
		if err == nil {
			*d = bv.(bool)
		}
		return err
	case *interface{}:
		*d = src
		return nil
	}
	dpv := reflect.ValueOf(dest)
	if dpv.Kind() != reflect.Ptr {
		return errors.New("destination not a pointer")
	}
	if dpv.IsNil() {
		return errNilPtr
	}

	if !sv.IsValid() {
		sv = reflect.ValueOf(src)
	}

	dv := reflect.Indirect(dpv)
	if dv.Kind() == sv.Kind() {
		dv.Set(sv)
		return nil
	}

	switch dv.Kind() {
	case reflect.Ptr:
		if src == nil {
			dv.Set(reflect.Zero(dv.Type()))
			return nil
		} else {
			dv.Set(reflect.New(dv.Type().Elem()))
			return i.convert(dv.Interface(), src)
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		s := i.asString(src)
		i64, err := strconv.ParseInt(s, 10, dv.Type().Bits())
		if err != nil {
			return fmt.Errorf("converting string %q to a %s: %v", s, dv.Kind(), err)
		}
		dv.SetInt(i64)
		return nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		s := i.asString(src)
		u64, err := strconv.ParseUint(s, 10, dv.Type().Bits())
		if err != nil {
			return fmt.Errorf("converting string %q to a %s: %v", s, dv.Kind(), err)
		}
		dv.SetUint(u64)
		return nil
	case reflect.Float32, reflect.Float64:
		s := i.asString(src)
		f64, err := strconv.ParseFloat(s, dv.Type().Bits())
		if err != nil {
			return fmt.Errorf("converting string %q to a %s: %v", s, dv.Kind(), err)
		}
		dv.SetFloat(f64)
		return nil
	}

	return fmt.Errorf("unsupported driver -> Scan pair: %T -> %T", src, dest)
}

func (i *Item) cloneBytes(b []byte) []byte {
	if b == nil {
		return nil
	} else {
		c := make([]byte, len(b))
		copy(c, b)
		return c
	}
}

func (i *Item) asString(src interface{}) string {
	switch v := src.(type) {
	case string:
		return v
	case []byte:
		return string(v)
	}
	rv := reflect.ValueOf(src)
	switch rv.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(rv.Int(), 10)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return strconv.FormatUint(rv.Uint(), 10)
	case reflect.Float64:
		return strconv.FormatFloat(rv.Float(), 'g', -1, 64)
	case reflect.Float32:
		return strconv.FormatFloat(rv.Float(), 'g', -1, 32)
	case reflect.Bool:
		return strconv.FormatBool(rv.Bool())
	}
	return fmt.Sprintf("%v", src)
}

func (i *Item) asBytes(buf []byte, rv reflect.Value) (b []byte, ok bool) {
	switch rv.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.AppendInt(buf, rv.Int(), 10), true
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return strconv.AppendUint(buf, rv.Uint(), 10), true
	case reflect.Float32:
		return strconv.AppendFloat(buf, rv.Float(), 'g', -1, 32), true
	case reflect.Float64:
		return strconv.AppendFloat(buf, rv.Float(), 'g', -1, 64), true
	case reflect.Bool:
		return strconv.AppendBool(buf, rv.Bool()), true
	case reflect.String:
		s := rv.String()
		return append(buf, s...), true
	}
	return
}
