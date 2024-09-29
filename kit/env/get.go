package env

import "os"

func GetStringOrDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func GetBoolOrDefault(key string, defaultValue bool) bool {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value == "true"
}

func LookupString(key string) (string, bool) {
	value, found := os.LookupEnv(key)
	return value, found
}
