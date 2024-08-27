package config

import "os"

func IsVerbose() bool {
	return os.Getenv("VERBOSE") == "true"
}
