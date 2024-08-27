//go:build print
// +build print

package platform

import (
	"encoding/json"
	"fmt"

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
	var prepared []string
	for _, word := range words {
		output, err := json.MarshalIndent(word, "  ", " ")
		if err != nil {
			return err
		}
		prepared = append(prepared, string(output))
	}
	err := PartitionAndPersist(bookName, prepared)
	if err != nil {
		return err
	}
	return nil
}

func PrepareAndPersistWlc(tableName string, bookName string, words []models.WlcWord) error {
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

func PartitionAndPersist(bookName string, prepared []string) error {
	PartitionSize := getPartitionSize()
	fmt.Printf("Partition size: %d\n", PartitionSize)
	segmentNumber := 1
	fmt.Printf("Saving %s (%d words)...\n", bookName, len(prepared))
	for idxRange := range util.Partition(len(prepared), PartitionSize) {
		// fmt.Printf("Partition: %d %d %d\n", idxRange.Low, idxRange.High, idxRange.High-idxRange.Low)
		segment := prepared[idxRange.Low:idxRange.High]
		err := persist(segment)
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

func persist(words []string) error {
	for _, obj := range words {
		fmt.Printf("Length: %d\n", len(obj))
	}
	return nil
}

func PostPersistWLC(tableName string) error {
	return nil
}

func PostPersistGNT(tableName string) error {
	return nil
}
