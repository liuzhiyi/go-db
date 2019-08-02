package db

type Filter struct {
	bucket map[string][]map[string]interface{}
	fileds []string
}

func NewFilter() *Filter {
	f := new(Filter)
	f.bucket =  make(map[string][]map[string]interface{})

	return f
}

func (f *Filter) SetCondition(field, key string, value interface{}) {
	condition := make(map[string]interface{})
	condition[key] = value
	if _, exist := f.bucket[field]; !exist {
		f.fileds = append(f.fileds, field)
	}
	f.bucket[field] = append(f.bucket[field], condition)
}
