package engine

func interfaceGet[T any](m map[string]interface{}, key string, defaultVal T) T {
	val, ok := m[key].(T)
	if !ok {
		return defaultVal
	}
	return val
}
