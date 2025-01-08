package main

import (
	"encoding/xml"
	"fmt"
	"os"
)

// Define structures for XML
type Lexicon struct {
	Entries []Entry `xml:"part>section>entry"`
}

type Entry struct {
	Word        string   `xml:"w"`
	Definitions []string `xml:"def"`
}

func main() {
	// Sample XML file (replace with your file path)
	file, err := os.Open("/home/dbetz/gh/morph/data/HebrewLexicon/BrownDriverBriggs.xml")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	// Decode the XML
	var lexicon Lexicon
	decoder := xml.NewDecoder(file)
	if err := decoder.Decode(&lexicon); err != nil {
		fmt.Println("Error decoding XML:", err)
		return
	}

	// Walk through entries and print words and definitions
	for _, entry := range lexicon.Entries {
		fmt.Printf("Word: %s\n", entry.Word)
		for _, def := range entry.Definitions {
			fmt.Printf("  Definition: %s\n", def)
		}
	}
}
