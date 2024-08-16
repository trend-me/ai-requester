package builders

func BuildMetadata(m1, m2 map[string]interface{}) map[string]interface{} {
	m := make(map[string]interface{})
	for k, v := range m1 {
		m[k] = v
	}
	for k, v := range m2 {
		m[k] = v
	}
	return m
}
