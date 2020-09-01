// +build gcp

package main

import (
	"context"
	"errors"
	"fmt"
	"os"

	"cloud.google.com/go/datastore"
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

type Saver func(context.Context, int, int, *datastore.Client) ([]*datastore.Key, error)

func getPartitionSize() int {
	return 200
}

func validateCloudConfig() error {
	projectID := os.Getenv("GCP_PROJECT_ID")
	if len(projectID) == 0 {
		return errors.New("GCP_PROJECT_ID is required.")
	}
	return nil
}

func prepareAndPersistWlc(tableName string, bookName string, words []wlcWord) error {
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
	return partitionAndPersist(tableName, bookName, len(prepared), f)
}

func prepareAndPersistGnt(tableName string, bookName string, words []gntWord) error {
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
	return partitionAndPersist(tableName, bookName, len(words), f)
}

func partitionAndPersist(tableName string, bookName string, size int, f Saver) error {
	partitionSize := getPartitionSize()
	fmt.Printf("Partition size: %d\n", partitionSize)
	segmentNumber := 1
	fmt.Printf("Saving %s (%d words)...\n", bookName, size)
	for idxRange := range Partition(size, partitionSize) {
		err := persist(idxRange.Low, idxRange.High, f)
		if err != nil {
			return err
		}
		percent := (float64(segmentNumber) * float64((partitionSize)) / float64(size)) * 100
		if percent > 100 {
			percent = 100
		}
		fmt.Printf("%s %0.2f%% complete\n", bookName, percent)
		segmentNumber++
	}
	return nil
}

func persist(start int, end int, f Saver) error {
	ctx := context.Background()
	projectID := os.Getenv("GCP_PROJECT_ID")
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

func postPersistWLC(tableName string) error {
	return nil
}

func postPersistGNT(tableName string) error {
	return nil
}
