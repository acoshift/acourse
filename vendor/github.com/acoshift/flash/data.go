package flash

// Data is the flash data
type Data map[string][]interface{}

// Set sets value to data
func (d Data) Set(key string, value interface{}) {
	d[key] = []interface{}{value}
}

// Add adds value to data
func (d Data) Add(key string, value interface{}) {
	d[key] = append(d[key], value)
}

// Get gets value from data
func (d Data) Get(key string) interface{} {
	if d == nil {
		return nil
	}
	if len(d[key]) == 0 {
		return nil
	}
	return d[key][0]
}

// GetString gets string from data
func (d Data) GetString(key string) string {
	r, _ := d.Get(key).(string)
	return r
}

// GetInt gets int from data
func (d Data) GetInt(key string) int {
	r, _ := d.Get(key).(int)
	return r
}

// GetInt64 gets int64 from data
func (d Data) GetInt64(key string) int64 {
	r, _ := d.Get(key).(int64)
	return r
}

// GetFloat32 gets float32 from data
func (d Data) GetFloat32(key string) float32 {
	r, _ := d.Get(key).(float32)
	return r
}

// GetFloat64 gets float64 from data
func (d Data) GetFloat64(key string) float64 {
	r, _ := d.Get(key).(float64)
	return r
}

// GetBool gets bool from data
func (d Data) GetBool(key string) bool {
	r, _ := d.Get(key).(bool)
	return r
}

// Del deletes key from data
func (d Data) Del(key string) {
	d[key] = nil
}

// Has checks is flash has a given key
func (d Data) Has(key string) bool {
	if d == nil {
		return false
	}
	return len(d[key]) > 0
}
