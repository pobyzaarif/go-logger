package database

// DatabaseLog : database log wrapper
func DatabaseLog(data map[string]interface{}) map[string]interface{} {
	return map[string]interface{}{
		"__gologger__": 1,
		"db":           data,
	}
}
