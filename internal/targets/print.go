package targets

import (
	"encoding/json"
	"fmt"

	"github.com/davidbetz/morph/internal/models"
	"github.com/davidbetz/morph/internal/util"
)

type Print struct{}

func CreatePrint() *Print {
	return &Print{}
}

func (t *Print) getPartitionSize() int {
	return 100
}

func (t *Print) ValidateCloudConfig() error {
	return nil
}

func (t *Print) unifiedPersist(tableName string, bookName string, words []interface{}) error {
	var prepared []string
	for _, word := range words {
		output, err := json.MarshalIndent(word, "  ", " ")
		if err != nil {
			return err
		}
		prepared = append(prepared, string(output))
	}
	err := t.PartitionAndPersist(bookName, prepared)
	if err != nil {
		return err
	}
	return nil
}

func (t *Print) PrepareAndPersistWlc(tableName string, bookName string, words []models.WlcWord) error {
	var taco []interface{}
	m, _ := json.Marshal(words)
	err := json.Unmarshal(m, &taco)
	if err != nil {
		return err
	}
	return t.unifiedPersist(tableName, bookName, taco)
}

func (t *Print) PrepareAndPersistGnt(tableName string, bookName string, words []models.GntWord) error {
	var taco []interface{}
	m, _ := json.Marshal(words)
	err := json.Unmarshal(m, &taco)
	if err != nil {
		return err
	}
	return t.unifiedPersist(tableName, bookName, taco)
}

func (t *Print) PartitionAndPersist(bookName string, prepared []string) error {
	PartitionSize := t.getPartitionSize()
	fmt.Printf("Partition size: %d\n", PartitionSize)
	segmentNumber := 1
	fmt.Printf("Saving %s (%d words)...\n", bookName, len(prepared))
	for idxRange := range util.Partition(len(prepared), PartitionSize) {
		// fmt.Printf("Partition: %d %d %d\n", idxRange.Low, idxRange.High, idxRange.High-idxRange.Low)
		segment := prepared[idxRange.Low:idxRange.High]
		err := t.persist(segment)
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

func (t *Print) persist(words []string) error {
	for _, obj := range words {
		fmt.Printf("Length: %d\n", len(obj))
	}
	return nil
}

func (t *Print) PostPersistWLC(tableName string) error {
	return nil
}

func (t *Print) PostPersistGNT(tableName string) error {
	return nil
}

func (t *Print) PersistCounts(bookName string, buckets map[int]map[string]models.WordCount) error {
	return nil
}

func (t *Print) PersistRender(bookName string, words []string) error {
	return nil
}

func (t *Print) PersistStrongs(PersistStrongs string, data []models.Lemma) error {
	return models.NewNotImplementedError("PersistStrongs")
}

func (t *Print) PersistMacula(language string, data []models.MaculaWord) error {
	return models.NewNotImplementedError("PersistMacula")
}
