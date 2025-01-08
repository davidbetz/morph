package targets

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

type Aws struct {
}

func CreateAws() *Aws {
	return &Aws{}
}

func (t *Aws) getPartitionSize() int {
	return 25
}

func (t *Aws) createSession() (*session.Session, error) {
	session, err := session.NewSession()
	if err != nil {
		util.Errorf(err.Error())
	}
	return session, err
}

func (t *Aws) createAttributeValue(word interface{}) (map[string]*dynamodb.AttributeValue, error) {
	av, err := dynamodbattribute.MarshalMap(word)
	if err != nil {
		return nil, fmt.Errorf("MarshalMap error %s", err.Error())
	}
	return av, nil
}

func (t *Aws) unifiedPersist(tableName string, bookName string, words []interface{}) error {
	prepared := make([]*dynamodb.WriteRequest, len(words))
	for i, word := range words {
		av, err := t.createAttributeValue(word)
		if err != nil {
			return err
		}
		prepared[i] = &dynamodb.WriteRequest{
			PutRequest: &dynamodb.PutRequest{
				Item: av,
			},
		}
	}
	return t.PartitionAndPersist(tableName, bookName, prepared)
}

func (t *Aws) PrepareAndPersistWlc(tableName string, bookName string, words []models.WlcWord) error {
	var taco []interface{}
	m, _ := json.Marshal(words)
	err := json.Unmarshal(m, &taco)
	if err != nil {
		return err
	}
	return t.unifiedPersist(tableName, bookName, taco)
}

func (t *Aws) PrepareAndPersistGnt(tableName string, bookName string, words []models.GntWord) error {
	var taco []interface{}
	m, _ := json.Marshal(words)
	err := json.Unmarshal(m, &taco)
	if err != nil {
		return err
	}
	return t.unifiedPersist(tableName, bookName, taco)
}

func (t *Aws) PartitionAndPersist(tableName string, bookName string, prepared []*dynamodb.WriteRequest) error {
	PartitionSize := t.getPartitionSize()
	fmt.Printf("Partition size: %d\n", PartitionSize)
	segmentNumber := 1
	fmt.Printf("Saving %s (%d words)...\n", bookName, len(prepared))
	for idxRange := range util.Partition(len(prepared), PartitionSize) {
		// fmt.Printf("Partition: %d %d %d\n", idxRange.Low, idxRange.High, idxRange.High-idxRange.Low)
		segment := prepared[idxRange.Low:idxRange.High]
		err := t.persist(tableName, segment)
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

func (t *Aws) ValidateCloudConfig() error {
	return nil
}

func (t *Aws) persist(tableName string, items []*dynamodb.WriteRequest) error {
	sess, err := t.createSession()
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

func (t *Aws) PostPersistWLC(tableName string) error {
	return nil
}

func (t *Aws) PostPersistGNT(tableName string) error {
	return nil
}

func (t *Aws) PersistCounts(bookName string, buckets map[int]map[string]models.WordCount) error {
	return nil
}

func (t *Aws) PersistRender(bookName string, words []string) error {
	return nil
}

func (t *Aws) PersistStrongs(PersistStrongs string, data []models.Lemma) error {
	return models.NewNotImplementedError("PersistStrongs")
}

func (t *Aws) PersistMacula(language string, data []models.MaculaWord) error {
	return models.NewNotImplementedError("PersistMacula")
}
