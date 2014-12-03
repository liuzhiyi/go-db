package data

import (
	"encoding/json"
	"strconv"
)

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
