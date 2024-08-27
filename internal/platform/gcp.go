//go:build gcp
// +build gcp

package platform

import (
	"context"
	"errors"
	"fmt"
	"os"

	"cloud.google.com/go/datastore"
	"github.com/davidbetz/morph/internal/models"
	"github.com/davidbetz/morph/internal/util"
)

type wlcWordDataStoreEntity struct {
	Codes      string `datastore:"codes"`
	Language   string `datastore:"language"`
	Lemma      string `datastore:"lemma"`
	ID         string `datastore:"coreid"`
	Morphology string `datastore:"morphology"`
	SequenceID int64  `datastore:"id"`
	Verse      string `datastore:"verse"`
}

type saver func(context.Context, int, int, *datastore.Client) ([]*datastore.Key, error)

func getPartitionSize() int {
	return 200
}

func ValidateCloudConfig() error {
	projectID := os.Getenv("PROJECT_ID")
	if len(projectID) == 0 {
		return errors.New("PROJECT_ID is required.")
	}
	return nil
}

func PrepareAndPersistWlc(tableName string, bookName string, words []models.WlcWord) error {
	var keys []*datastore.Key
	var prepared []wlcWordDataStoreEntity
	for _, word := range words {
		keys = append(keys, datastore.NameKey(tableName, fmt.Sprintf("%d", word.SequenceID), nil))
		prepared = append(prepared, wlcWordDataStoreEntity{
			Codes:      word.Codes,
			Language:   word.Language,
			Lemma:      word.Lemma,
			ID:         word.ID,
			Morphology: word.MorphologyString,
			SequenceID: word.SequenceID,
			Verse:      word.Verse,
		})
	}
	//+ strategy pattern bc of different types
	f := func(ctx context.Context, start int, end int, client *datastore.Client) ([]*datastore.Key, error) {
		results, err := client.PutMulti(ctx, keys[start:end], prepared[start:end])
		if err != nil {
			return nil, err
		}
		return results, nil
	}
	return PartitionAndPersist(tableName, bookName, len(prepared), f)
}

func PrepareAndPersistGnt(tableName string, bookName string, words []models.GntWord) error {
	var keys []*datastore.Key
	for _, key := range words {
		keys = append(keys, datastore.NameKey(tableName, fmt.Sprintf("%d", key.ID), nil))
	}
	f := func(ctx context.Context, start int, end int, client *datastore.Client) ([]*datastore.Key, error) {
		results, err := client.PutMulti(ctx, keys[start:end], words[start:end])
		if err != nil {
			return nil, err
		}
		return results, nil
	}
	return PartitionAndPersist(tableName, bookName, len(words), f)
}

func PartitionAndPersist(tableName string, bookName string, size int, f saver) error {
	PartitionSize := getPartitionSize()
	fmt.Printf("Partition size: %d\n", PartitionSize)
	segmentNumber := 1
	fmt.Printf("Saving %s (%d words)...\n", bookName, size)
	for idxRange := range util.Partition(size, PartitionSize) {
		err := persist(idxRange.Low, idxRange.High, f)
		if err != nil {
			return err
		}
		percent := (float64(segmentNumber) * float64((PartitionSize)) / float64(size)) * 100
		if percent > 100 {
			percent = 100
		}
		fmt.Printf("%s %0.2f%% complete\n", bookName, percent)
		segmentNumber++
	}
	return nil
}

func persist(start int, end int, f saver) error {
	ctx := context.Background()
	projectID := os.Getenv("PROJECT_ID")
	client, err := datastore.NewClient(ctx, projectID)
	if err != nil {
		return err
	}
	if f == nil {
		return errors.New("f is nil")
	}
	_, err = f(ctx, start, end, client)
	if err != nil {
		return err
	}
	return nil
}

func PostPersistWLC(tableName string) error {
	return nil
}

func PostPersistGNT(tableName string) error {
	return nil
}
