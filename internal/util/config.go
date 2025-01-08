package util

import "os"

func LoadEnvVarDef(name string, def string) string {
	value := os.Getenv(name)
	if value == "" {
		return def
	}
	return value
}
