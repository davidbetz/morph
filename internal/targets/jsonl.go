package targets

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"sort"

	"github.com/davidbetz/morph/internal/models"
	"github.com/davidbetz/morph/internal/util"
)

type JsonL struct {
}

func CreateJsonl() *JsonL {
	return &JsonL{}
}

func (t *JsonL) getPartitionSize() int {
	return 100
}

func (t *JsonL) ValidateCloudConfig() error {
	return nil
}

func (t *JsonL) unifiedPersist(tableName string, bookName string, words []interface{}) error {
	prepared := make([][]byte, len(words))
	for _, word := range words {
		output, err := json.Marshal(word)
		if err != nil {
			return err
		}
		output = append(output, byte('\n'))
		prepared = append(prepared, output)
	}
	err := t.PartitionAndPersist(tableName, bookName, prepared)
	if err != nil {
		return err
	}
	return nil
}

func (t *JsonL) PrepareAndPersistWlc(tableName string, bookName string, words []models.WlcWord) error {
	//+ trick to unify the logic; fine when perf isn't an issue
	var taco []interface{}
	m, _ := json.Marshal(words)
	err := json.Unmarshal(m, &taco)
	if err != nil {
		return err
	}
	return t.unifiedPersist(tableName, bookName, taco)
}

func (t *JsonL) PrepareAndPersistGnt(tableName string, bookName string, words []models.GntWord) error {
	var taco []interface{}
	m, _ := json.Marshal(words)
	err := json.Unmarshal(m, &taco)
	if err != nil {
		return err
	}
	return t.unifiedPersist(tableName, bookName, taco)
}

func (t *JsonL) PartitionAndPersist(tableName string, bookName string, prepared [][]byte) error {
	PartitionSize := t.getPartitionSize()
	fmt.Printf("Partition size: %d\n", PartitionSize)
	segmentNumber := 1
	fmt.Printf("Saving %s (%d words)...\n", bookName, len(prepared))
	for idxRange := range util.Partition(len(prepared), PartitionSize) {
		// fmt.Printf("Partition: %d %d %d\n", idxRange.Low, idxRange.High, idxRange.High-idxRange.Low)
		segment := prepared[idxRange.Low:idxRange.High]
		err := t.persist(tableName, bookName, segment)
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

func (t *JsonL) persist(tableName string, bookName string, words [][]byte) error {
	folder := path.Join("./output", tableName)
	if _, err := os.Stat(folder); os.IsNotExist(err) {
		_ = os.MkdirAll(folder, 0777)
	}
	filename := path.Join(folder, bookName) + ".jsonl"
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}
	defer func() {
		if cerr := f.Close(); cerr != nil && err == nil {
			err = cerr
		}
	}()
	for _, word := range words {
		if _, err = f.Write(word); err != nil {
			return err
		}
	}
	return f.Sync()
}

func (t *JsonL) PostPersistWLC(tableName string) error {
	return nil
}

func (t *JsonL) PostPersistGNT(tableName string) error {
	return nil
}

func (t *JsonL) PersistStrongs(language string, words []models.Lemma) error {
	newline := []byte("\n")
	folder := path.Join("./output", "strongs")
	if _, err := os.Stat(folder); os.IsNotExist(err) {
		_ = os.MkdirAll(folder, 0777)
	}
	filename := path.Join(folder, "strongs-"+language) + ".jsonl"
	if _, err := os.Stat(filename); err == nil {
		if err := os.Remove(filename); err != nil {
			return err
		}
	}
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}
	defer func() {
		if cerr := f.Close(); cerr != nil && err == nil {
			err = cerr
		}
	}()
	for _, word := range words {
		jsonData, err := json.Marshal(word)
		if err != nil {
			return err
		}
		jsonData = append(jsonData, newline...)
		if _, err = f.Write(jsonData); err != nil {
			return err
		}
	}
	return f.Sync()
}

func (t *JsonL) PersistMacula(language string, data []models.MaculaWord) error {
	newline := []byte("\n")
	folder := path.Join("./output", "macula")
	if _, err := os.Stat(folder); os.IsNotExist(err) {
		_ = os.MkdirAll(folder, 0777)
	}
	filename := path.Join(folder, "macula-"+language) + ".jsonl"
	if _, err := os.Stat(filename); err == nil {
		if err := os.Remove(filename); err != nil {
			return err
		}
	}
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}
	defer func() {
		if cerr := f.Close(); cerr != nil && err == nil {
			err = cerr
		}
	}()
	for _, word := range data {
		jsonData, err := json.Marshal(word)
		if err != nil {
			return err
		}
		jsonData = append(jsonData, newline...)
		if _, err = f.Write(jsonData); err != nil {
			return err
		}
	}
	return f.Sync()
}

func (t *JsonL) PersistCounts(bookName string, buckets map[int]map[string]models.WordCount) error {
	newline := []byte("\n")
	folder := path.Join("./output/counts", bookName)
	if _, err := os.Stat(folder); os.IsNotExist(err) {
		_ = os.MkdirAll(folder, 0777)
	}
	for bucketName, data := range buckets {
		filename := path.Join(folder, fmt.Sprintf("%d", bucketName)) + ".jsonl"
		// fmt.Printf("Persisting %s\n", filename)
		f, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
		if err != nil {
			return err
		}
		defer func() {
			if cerr := f.Close(); cerr != nil && err == nil {
				err = cerr
			}
		}()
		sortedWords := make([]models.WordCount, 0, len(data))
		for _, wordCount := range data {
			sortedWords = append(sortedWords, wordCount)
		}
		sort.Slice(sortedWords, func(i, j int) bool {
			return sortedWords[i].Count > sortedWords[j].Count
		})
		for _, wordCount := range sortedWords {
			jsonData, err := json.Marshal(wordCount)
			if err != nil {
				return err
			}
			jsonData = append(jsonData, newline...)
			if _, err = f.Write(jsonData); err != nil {
				return err
			}
		}
		err = f.Sync()
		if err != nil {
			return nil
		}
	}
	return nil
}

func (t *JsonL) PersistRender(bookName string, words []string) error {
	return models.NewNotImplementedError("PersistRender")
}
