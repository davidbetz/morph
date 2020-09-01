package main

import (
	"encoding/json"
	"fmt"
)

func debug(text string) {
	if verbose {
		fmt.Print(text)
	}
}

func dump(name string, obj interface{}) {
	fmt.Println("DUMP")
	b, err := json.MarshalIndent(obj, "", "  ")
	if err != nil {
		fmt.Println("JSON ERROR:", err)
	}
	fmt.Printf(">>>>%s\n%s\n<<<<%s\n", name, string(b), name)
}
