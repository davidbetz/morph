//go:build json
// +build json

package platform

import (
	"encoding/json"
	"fmt"
	"os"
	"path"

	"github.com/davidbetz/morph/internal/models"
	"github.com/davidbetz/morph/internal/util"
)

func getPartitionSize() int {
	return 100
}

func ValidateCloudConfig() error {
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
	err := PartitionAndPersist(tableName, bookName, prepared)
	if err != nil {
		return err
	}
	return nil
}

func PrepareAndPersistWlc(tableName string, bookName string, words []models.WlcWord) error {
	//+ trick to unify the logic; fine when perf isn't an issue
	var taco []interface{}
	m, _ := json.Marshal(words)
	json.Unmarshal(m, &taco)
	return unifiedPersist(tableName, bookName, taco)
}

func PrepareAndPersistGnt(tableName string, bookName string, words []models.GntWord) error {
	var taco []interface{}
	m, _ := json.Marshal(words)
	json.Unmarshal(m, &taco)
	return unifiedPersist(tableName, bookName, taco)
}

func PartitionAndPersist(tableName string, bookName string, prepared [][]byte) error {
	PartitionSize := getPartitionSize()
	fmt.Printf("Partition size: %d\n", PartitionSize)
	segmentNumber := 1
	fmt.Printf("Saving %s (%d words)...\n", bookName, len(prepared))
	for idxRange := range util.Partition(len(prepared), PartitionSize) {
		// fmt.Printf("Partition: %d %d %d\n", idxRange.Low, idxRange.High, idxRange.High-idxRange.Low)
		segment := prepared[idxRange.Low:idxRange.High]
		err := persist(tableName, bookName, segment)
		if err != nil {
			return err
		}
		percent := (float64(segmentNumber) * float64((PartitionSize)) / float64(len(prepared))) * 100
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

func PostPersistWLC(tableName string) error {
	return nil
}

func PostPersistGNT(tableName string) error {
	return nil
}
