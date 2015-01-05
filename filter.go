package db

type Filter map[string][]map[string]interface{}

func NewFilter() Filter {
	f := make(Filter)
	return f
}

func (f Filter) SetCondition(field, key string, value interface{}) {
	condition := make(map[string]interface{})
	condition[key] = value
	f[field] = append(f[field], condition)
}
