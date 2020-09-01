// +build json

package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
)

func getPartitionSize() int {
	return 100
}

func validateCloudConfig() error {
	return nil
}

func unifiedPersist(tableName string, bookName string, words []interface{}) error {
	prepared := make([][]byte, len(words))
	for _, word := range words {
		output, err := json.Marshal(word)
		if err != nil {
			return err
		}
		output = append(output, byte('\n'))
		prepared = append(prepared, output)
	}
	err := partitionAndPersist(tableName, bookName, prepared)
	if err != nil {
		return err
	}
	return nil
}

func prepareAndPersistWlc(tableName string, bookName string, words []wlcWord) error {
	//+ trick to unify the logic; fine when perf isn't an issue
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

func partitionAndPersist(tableName string, bookName string, prepared [][]byte) error {
	partitionSize := getPartitionSize()
	fmt.Printf("Partition size: %d\n", partitionSize)
	segmentNumber := 1
	fmt.Printf("Saving %s (%d words)...\n", bookName, len(prepared))
	for idxRange := range Partition(len(prepared), partitionSize) {
		// fmt.Printf("Partition: %d %d %d\n", idxRange.Low, idxRange.High, idxRange.High-idxRange.Low)
		segment := prepared[idxRange.Low:idxRange.High]
		err := persist(tableName, bookName, segment)
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

func persist(tableName string, bookName string, words [][]byte) error {
	folder := path.Join("./output", tableName)
	if _, err := os.Stat(folder); os.IsNotExist(err) {
		os.MkdirAll(folder, 0777)
	}
	filename := path.Join(folder, bookName) + ".jsonl"
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}
	defer f.Close()
	for _, word := range words {
		if _, err = f.Write(word); err != nil {
			return err
		}
	}
	f.Sync()
	return nil
}

func postPersistWLC(tableName string) error {
	return nil
}

func postPersistGNT(tableName string) error {
	return nil
}
