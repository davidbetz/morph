//go:build aws
// +build aws

package platform

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/davidbetz/morph/internal/models"
	"github.com/davidbetz/morph/internal/util"
)

func getPartitionSize() int {
	return 25
}

func createSession() (*session.Session, error) {
	session, err := session.NewSession()
	if err != nil {
		util.Errorf(err.Error())
	}
	return session, err
}

func createAttributeValue(word interface{}) (map[string]*dynamodb.AttributeValue, error) {
	av, err := dynamodbattribute.MarshalMap(word)
	if err != nil {
		return nil, fmt.Errorf("MarshalMap error %s", err.Error())
	}
	return av, nil
}

func unifiedPersist(tableName string, bookName string, words []interface{}) error {
	prepared := make([]*dynamodb.WriteRequest, len(words))
	for i, word := range words {
		av, err := createAttributeValue(word)
		if err != nil {
			return err
		}
		prepared[i] = &dynamodb.WriteRequest{
			PutRequest: &dynamodb.PutRequest{
				Item: av,
			},
		}
	}
	return PartitionAndPersist(tableName, bookName, prepared)
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

func PartitionAndPersist(tableName string, bookName string, prepared []*dynamodb.WriteRequest) error {
	PartitionSize := getPartitionSize()
	fmt.Printf("Partition size: %d\n", PartitionSize)
	segmentNumber := 1
	fmt.Printf("Saving %s (%d words)...\n", bookName, len(prepared))
	for idxRange := range util.Partition(len(prepared), PartitionSize) {
		// fmt.Printf("Partition: %d %d %d\n", idxRange.Low, idxRange.High, idxRange.High-idxRange.Low)
		segment := prepared[idxRange.Low:idxRange.High]
		err := persist(tableName, segment)
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

func ValidateCloudConfig() error {
	return nil
}

func persist(tableName string, items []*dynamodb.WriteRequest) error {
	sess, err := createSession()
	if err != nil {
		return fmt.Errorf("NewSession error %s", err.Error())
	}

	records := make(map[string][]*dynamodb.WriteRequest, 1)
	notdone := true
	retry := 0
	backoff := 1
	for notdone {
		records[tableName] = items
		input := &dynamodb.BatchWriteItemInput{
			RequestItems: records,
		}
		svc := dynamodb.New(sess)
		response, err := svc.BatchWriteItem(input)
		if err != nil {
			return err
		}
		items = response.UnprocessedItems[tableName]
		if len(items) == 0 {
			notdone = false
			continue
		}
		retry++
		time.Sleep(time.Duration(backoff) * time.Second)
		fmt.Printf("DYNAMODB BACKING OFF (%d) | Left to process: %d | Backoff: %ds\n", retry, len(items), backoff)
		backoff *= 2
	}
	return nil
}

func PostPersistWLC(tableName string) error {
	return nil
}

func PostPersistGNT(tableName string) error {
	return nil
}
