// +build print

package main

import (
	"encoding/json"
	"fmt"
)

func getPartitionSize() int {
	return 100
}

func validateCloudConfig() error {
	return nil
}

func unifiedPersist(tableName string, bookName string, words []interface{}) error {
	var prepared []string
	for _, word := range words {
		output, err := json.MarshalIndent(word, "  ", " ")
		if err != nil {
			return err
		}
		prepared = append(prepared, string(output))
	}
	err := partitionAndPersist(bookName, prepared)
	if err != nil {
		return err
	}
	return nil
}

func prepareAndPersistWlc(tableName string, bookName string, words []wlcWord) error {
	var taco []interface{}
	m, _ := json.Marshal(words)
	json.Unmarshal(m, &taco)
	return unifiedPersist(tableName, bookName, taco)
}

func prepareAndPersistGnt(tableName string, bookName string, words []gntWord) error {
	var taco []interface{}
	m, _ := json.Marshal(words)
	json.Unmarshal(m, &taco)
	return unifiedPersist(tableName, bookName, taco)
}

func partitionAndPersist(bookName string, prepared []string) error {
	partitionSize := getPartitionSize()
	fmt.Printf("Partition size: %d\n", partitionSize)
	segmentNumber := 1
	fmt.Printf("Saving %s (%d words)...\n", bookName, len(prepared))
	for idxRange := range Partition(len(prepared), partitionSize) {
		// fmt.Printf("Partition: %d %d %d\n", idxRange.Low, idxRange.High, idxRange.High-idxRange.Low)
		segment := prepared[idxRange.Low:idxRange.High]
		err := persist(segment)
		if err != nil {
			return err
		}
		percent := (float64(segmentNumber) * float64((partitionSize)) / float64(len(prepared))) * 100
		if percent > 100 {
			percent = 100
		}
		fmt.Printf("%s %0.2f%% complete\n", bookName, percent)
		segmentNumber++
	}
	return nil
}

func persist(words []string) error {
	for _, obj := range words {
		fmt.Printf("Length: %d\n", len(obj))
	}
	return nil
}

func postPersistWLC(tableName string) error {
	return nil
}

func postPersistGNT(tableName string) error {
	return nil
}
