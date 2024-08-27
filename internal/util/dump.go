package util

import (
	"encoding/json"
	"fmt"

	"github.com/davidbetz/morph/internal/config"
)

func Debug(text string) {
	if config.IsVerbose() {
		fmt.Print(text)
	}
}

func Dump(name string, obj interface{}) {
	fmt.Println("DUMP")
	b, err := json.MarshalIndent(obj, "", "  ")
	if err != nil {
		fmt.Println("JSON ERROR:", err)
	}
	fmt.Printf(">>>>%s\n%s\n<<<<%s\n", name, string(b), name)
}
