// +build azure

package main

//+ https://github.com/Azure/azure-sdk-for-go/blob/77258e94d84ea36012a72c0e0a1e2faa409c6396/storage/entity_test.go

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/Azure/azure-sdk-for-go/storage"
)

type Word struct {
	PartitionKey string
	RowKey       string
	Properties   map[string]interface{}
}

func getTableReference(tableName string) *storage.Table {
	sas := os.Getenv("AZURE_CS")
	client, err := storage.NewClientFromConnectionString(sas)
	if err != nil {
		log.Fatal(err)
	}
	tableService := client.GetTableService()
	return tableService.GetTableReference(tableName)
}

func getPartitionSize() int {
	return 1000
}

func validateCloudConfig() error {
	sas := os.Getenv("AZURE_CS")
	if len(sas) == 0 {
		return errors.New("AZURE_CS is required.")
	}
	return nil
}

func prepareAndPersistWlc(tableName string, bookName string, words []wlcWord) error {
	var prepared []Word
	for _, word := range words {
		preparedProperties := map[string]interface{}{
			"Lemma":  word.Lemma,
			"CoreID": word.ID,
			//+ separating each part to a different column creates far too many
			"MorphCodes": word.MorphologyString,
			"UniqueID":   word.Verse,
			"Codes":      word.Codes,
		}
		prepared = append(prepared, Word{
			PartitionKey: word.Verse,
			RowKey:       fmt.Sprintf("%d", word.SequenceID),
			Properties:   preparedProperties,
		})
	}
	return partitionAndPersist(tableName, bookName, prepared)
}

func prepareAndPersistGnt(tableName string, bookName string, words []gntWord) error {
	var prepared []Word
	for _, word := range words {
		prepared = append(prepared, Word{
			PartitionKey: word.Verse,
			RowKey:       fmt.Sprintf("%d", word.ID),
			Properties: map[string]interface{}{
				"Part":       word.Morphology.Part,
				"Person":     word.Morphology.Person,
				"Tense":      word.Morphology.Tense,
				"Voice":      word.Morphology.Voice,
				"Mood":       word.Morphology.Mood,
				"Case":       word.Morphology.Case,
				"Number":     word.Morphology.Number,
				"Gender":     word.Morphology.Gender,
				"Degree":     word.Morphology.Degree,
				"Text":       word.Text,
				"Word":       word.Word,
				"Normalized": word.Normalized,
				"Lemma":      word.Lemma,
				"Codes":      word.Codes,
			},
		})
	}
	return partitionAndPersist(tableName, bookName, prepared)
}

func partitionAndPersist(tableName string, bookName string, prepared []Word) error {
	partitionSize := getPartitionSize()
	fmt.Printf("Partition size: %d\n", partitionSize)
	segmentNumber := 1
	fmt.Printf("Saving %s (%d words)...\n", bookName, len(prepared))
	for idxRange := range Partition(len(prepared), partitionSize) {
		segment := prepared[idxRange.Low:idxRange.High]
		err := persist(tableName, segment)
		if err != nil {
			return err
		}
		percent := ((float64(segmentNumber) * float64((partitionSize)) / float64(len(prepared))) * 100)
		if percent > 100 {
			percent = 100
		}
		fmt.Printf("%s %0.2f%% complete\n", bookName, percent)
		segmentNumber++
	}
	return nil
}

func persist(tableName string, segment []Word) error {
	for _, word := range segment {
		table := getTableReference(tableName)
		entity := table.GetEntityReference(word.PartitionKey, word.RowKey)
		entity.Properties = word.Properties
		err := entity.InsertOrReplace(nil)
		if err != nil {
			return err
		}
	}
	return nil
}

func postPersistWLC(tableName string) error {
	return nil
}

func postPersistGNT(tableName string) error {
	return nil
}
