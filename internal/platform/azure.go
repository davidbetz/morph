//go:build azure
// +build azure

package platform

//+ https://github.com/Azure/azure-sdk-for-go/blob/77258e94d84ea36012a72c0e0a1e2faa409c6396/storage/entity_test.go

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/Azure/azure-sdk-for-go/storage"
	"github.com/davidbetz/morph/internal/models"
	"github.com/davidbetz/morph/internal/util"
)

type azureWord struct {
	PartitionKey string
	RowKey       string
	Properties   map[string]interface{}
}

func getTableReference(tableName string) *storage.Table {
	cs := os.Getenv("CS")
	client, err := storage.NewClientFromConnectionString(cs)
	if err != nil {
		log.Fatal(err)
	}
	tableService := client.GetTableService()
	return tableService.GetTableReference(tableName)
}

func getPartitionSize() int {
	return 1000
}

func ValidateCloudConfig() error {
	cs := os.Getenv("CS")
	if len(cs) == 0 {
		return errors.New("CS is required.")
	}
	return nil
}

func PrepareAndPersistWlc(tableName string, bookName string, words []models.WlcWord) error {
	var prepared []azureWord
	for _, word := range words {
		preparedProperties := map[string]interface{}{
			"Lemma":  word.Lemma,
			"CoreID": word.ID,
			//+ separating each part to a different column creates far too many
			"MorphCodes": word.MorphologyString,
			"UniqueID":   word.Verse,
			"Codes":      word.Codes,
		}
		prepared = append(prepared, azureWord{
			PartitionKey: word.Verse,
			RowKey:       fmt.Sprintf("%d", word.SequenceID),
			Properties:   preparedProperties,
		})
	}
	return PartitionAndPersist(tableName, bookName, prepared)
}

func PrepareAndPersistGnt(tableName string, bookName string, words []models.GntWord) error {
	var prepared []azureWord
	for _, word := range words {
		prepared = append(prepared, azureWord{
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
	return PartitionAndPersist(tableName, bookName, prepared)
}

func PartitionAndPersist(tableName string, bookName string, prepared []azureWord) error {
	PartitionSize := getPartitionSize()
	fmt.Printf("Partition size: %d\n", PartitionSize)
	segmentNumber := 1
	fmt.Printf("Saving %s (%d words)...\n", bookName, len(prepared))
	for idxRange := range util.Partition(len(prepared), PartitionSize) {
		segment := prepared[idxRange.Low:idxRange.High]
		err := persist(tableName, segment)
		if err != nil {
			return err
		}
		percent := ((float64(segmentNumber) * float64((PartitionSize)) / float64(len(prepared))) * 100)
		if percent > 100 {
			percent = 100
		}
		fmt.Printf("%s %0.2f%% complete\n", bookName, percent)
		segmentNumber++
	}
	return nil
}

func persist(tableName string, segment []azureWord) error {
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

func PostPersistWLC(tableName string) error {
	return nil
}

func PostPersistGNT(tableName string) error {
	return nil
}
