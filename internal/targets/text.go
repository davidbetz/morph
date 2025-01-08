package targets

import (
	"fmt"
	"os"
	"path"
	"sort"
	"strings"

	"github.com/davidbetz/morph/internal/models"
)

type Text struct {
}

func CreateText() *Text {
	return &Text{}
}

func (t *Text) getPartitionSize() int {
	return 100
}

func (t *Text) ValidateCloudConfig() error {
	return nil
}

func (t *Text) PrepareAndPersistWlc(tableName string, bookName string, words []models.WlcWord) error {
	return models.NewNotImplementedError("PrepareAndPersistWlc")
}

func (t *Text) PrepareAndPersistGnt(tableName string, bookName string, words []models.GntWord) error {
	return models.NewNotImplementedError("PrepareAndPersistGnt")
}

func (t *Text) PartitionAndPersist(tableName string, bookName string, prepared [][]byte) error {
	return models.NewNotImplementedError("PartitionAndPersist")
}

func (t *Text) persist(tableName string, bookName string, words [][]byte) error {
	return nil
}
func (t *Text) PostPersistWLC(tableName string) error {
	return models.NewNotImplementedError("PostPersistWLC")
}

func (t *Text) PostPersistGNT(tableName string) error {
	return models.NewNotImplementedError("PostPersistGNT")
}

func (t *Text) PersistStrongs(language string, words []models.Lemma) error {
	folder := path.Join("./output")
	if _, err := os.Stat(folder); os.IsNotExist(err) {
		_ = os.MkdirAll(folder, 0777)
	}
	filename := path.Join(folder, "strongs-"+language) + ".txt"
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
		line := []byte(fmt.Sprintf("%d: %s, %s\n", word.Number, word.Word, word.Gloss))
		if _, err = f.Write(line); err != nil {
			return err
		}
	}
	return f.Sync()
}

func (t *Text) PersistCounts(bookName string, buckets map[int]map[string]models.WordCount) error {
	newline := []byte("\n")
	folder := path.Join("./output/counts", bookName)
	if _, err := os.Stat(folder); os.IsNotExist(err) {
		_ = os.MkdirAll(folder, 0777)
	}
	for bucketName, data := range buckets {
		filename := path.Join(folder, fmt.Sprintf("%d", bucketName)) + ".txt"
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
			line := []byte(wordCount.String())
			line = append(line, newline...)
			if _, err = f.Write(line); err != nil {
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

func (t *Text) PersistRender(bookName string, words []string) error {
	folder := path.Join("./output/render")
	if _, err := os.Stat(folder); os.IsNotExist(err) {
		_ = os.MkdirAll(folder, 0777)
	}
	filename := path.Join(folder, bookName) + ".txt"
	fmt.Printf("Persisting %s\n", filename)
	text := strings.Join(words, " ")
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}
	defer func() {
		if cerr := f.Close(); cerr != nil && err == nil {
			err = cerr
		}
	}()
	if _, err = f.Write([]byte(text)); err != nil {
		return err
	}
	return f.Sync()
}

func (t *Text) PersistMacula(language string, data []models.MaculaWord) error {
	return models.NewNotImplementedError("PersistMacula")
}
