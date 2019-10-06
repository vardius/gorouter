package context

type (
	// Param object to hold request parameter
	Param struct {
		Key   string
		Value string
	}
	// Params slice returned from request context
	Params []Param
)

// Value of the request parameter by name
func (p Params) Value(key string) string {
	for i := range p {
		if p[i].Key == key {
			return p[i].Value
		}
	}
	return ""
}

// Set key value pair at index
func (p Params) Set(index uint8, key string, value string) {
	p[index].Value = value
	p[index].Key = key
}
