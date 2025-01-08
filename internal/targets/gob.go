package targets

import (
	"encoding/gob"
	"os"
	"path"

	"github.com/davidbetz/morph/internal/models"
)

type Gob struct {
}

func CreateGob() *Gob {
	return &Gob{}
}

func (t *Gob) getPartitionSize() int {
	return 100
}

func (t *Gob) ValidateCloudConfig() error {
	return nil
}

func (t *Gob) PrepareAndPersistWlc(tableName string, bookName string, words []models.WlcWord) error {
	return models.NewNotImplementedError("PrepareAndPersistWlc")
}

func (t *Gob) PrepareAndPersistGnt(tableName string, bookName string, words []models.GntWord) error {
	return models.NewNotImplementedError("PrepareAndPersistGnt")
}

func (t *Gob) PartitionAndPersist(tableName string, bookName string, prepared [][]byte) error {
	return models.NewNotImplementedError("PartitionAndPersist")
}

func (t *Gob) persist(tableName string, bookName string, words [][]byte) error {
	return nil
}
func (t *Gob) PostPersistWLC(tableName string) error {
	return models.NewNotImplementedError("PostPersistWLC")
}

func (t *Gob) PostPersistGNT(tableName string) error {
	return models.NewNotImplementedError("PostPersistGNT")
}

func (t *Gob) PersistStrongs(language string, words []models.Lemma) error {
	folder := path.Join("./output")
	if _, err := os.Stat(folder); os.IsNotExist(err) {
		_ = os.MkdirAll(folder, 0777)
	}
	filename := path.Join(folder, "strongs-"+language) + ".gob"
	if _, err := os.Stat(filename); err == nil {
		if err := os.Remove(filename); err != nil {
			return err
		}
	}
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer func() {
		if cerr := f.Close(); cerr != nil && err == nil {
			err = cerr
		}
	}()
	encoder := gob.NewEncoder(f)
	if err := encoder.Encode(words); err != nil {
		return err
	}
	return f.Sync()
}

func (t *Gob) PersistMacula(language string, words []models.MaculaWord) error {
	folder := path.Join("./output")
	if _, err := os.Stat(folder); os.IsNotExist(err) {
		_ = os.MkdirAll(folder, 0777)
	}
	filename := path.Join(folder, "macula-"+language) + ".gob"
	if _, err := os.Stat(filename); err == nil {
		if err := os.Remove(filename); err != nil {
			return err
		}
	}
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer func() {
		if cerr := f.Close(); cerr != nil && err == nil {
			err = cerr
		}
	}()
	encoder := gob.NewEncoder(f)
	if err := encoder.Encode(words); err != nil {
		return err
	}
	return f.Sync()
}

func (t *Gob) PersistCounts(bookName string, buckets map[int]map[string]models.WordCount) error {
	return models.NewNotImplementedError("PersistCounts")
}

func (t *Gob) PersistRender(bookName string, words []string) error {
	return models.NewNotImplementedError("PersistRender")
}
