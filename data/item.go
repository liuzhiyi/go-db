package data

import (
	"encoding/json"
	"errors"
	"strconv"
)

var errNilPtr = errors.New("destination pointer is nil")

type Item struct {
	data map[string]interface{}
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

func (i *Item) GetInt(key string) int {
	val, _ := strconv.Atoi(i.data[key])
	return val
}

func (i *Item) GetString()

func (i *Item) convert(dst, src interface{}) error {
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
			*d = cloneBytes(s)
			return nil
		case *[]byte:
			if d == nil {
				return errNilPtr
			}
			*d = cloneBytes(s)
			return nil
		case *RawBytes:
			if d == nil {
				return errNilPtr
			}
			*d = s
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
		case *RawBytes:
			if d == nil {
				return errNilPtr
			}
			*d = nil
			return nil
		}
	}
}
