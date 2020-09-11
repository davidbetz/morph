// +build mssql

package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"

	_ "github.com/denisenkom/go-mssqldb"
	mssql "github.com/denisenkom/go-mssqldb"
)

const (
	createGNTTable = `IF EXISTS (SELECT 0
		FROM INFORMATION_SCHEMA.TABLES
		WHERE TABLE_SCHEMA = 'dbo' AND TABLE_NAME = '{{ TABLE_NAME }}')
		BEGIN
			PRINT '{{ TABLE_NAME }} exists.'
		END
		ELSE
		BEGIN
			PRINT 'Creating {{ TABLE_NAME }}...'
			SET QUOTED_IDENTIFIER ON;
		
			CREATE TABLE [dbo].{{ TABLE_NAME }}
			(
				WordID AS CONVERT(nvarchar(200), JSON_VALUE(Content, '$.id')),
				VerseID AS CONVERT(nvarchar(200), JSON_VALUE(Content, '$.verse')),
				Text AS CONVERT(nvarchar(200), JSON_VALUE(Content, '$.text')),
				mssqlWord AS CONVERT(nvarchar(200), JSON_VALUE(Content, '$.word')),
				Normalized AS CONVERT(nvarchar(200), JSON_VALUE(Content, '$.normalized')),
				Lemma AS CONVERT(nvarchar(200), JSON_VALUE(Content, '$.lemma')),
				Codes AS CONVERT(nvarchar(200), JSON_VALUE(Content, '$.codes')),
				Content [nvarchar](max) NOT NULL
			);
		
			CREATE CLUSTERED INDEX Index{{ TABLE_NAME }}WordID ON {{ TABLE_NAME }} (WordID);
			ALTER TABLE [dbo].{{ TABLE_NAME }} ADD CONSTRAINT {{ TABLE_NAME }}ContentJson CHECK (ISJSON(Content)=1);
		
		END`
	createGNTIndexes = `
	CREATE CLUSTERED INDEX Index{{ TABLE_NAME }}WordID ON {{ TABLE_NAME }} (WordID);
	ALTER TABLE [dbo].{{ TABLE_NAME }} ADD CONSTRAINT {{ TABLE_NAME }}ContentJson CHECK (ISJSON(Content)=1);
	`

	createWLCTable = `IF EXISTS (SELECT 0
			FROM INFORMATION_SCHEMA.TABLES
			WHERE TABLE_SCHEMA = 'dbo' AND TABLE_NAME = '{{ TABLE_NAME }}')
			BEGIN
				PRINT '{{ TABLE_NAME }} exists.'
			END
			ELSE
			BEGIN
				PRINT 'Creating {{ TABLE_NAME }}...'
				SET QUOTED_IDENTIFIER ON;
			
				CREATE TABLE [dbo].{{ TABLE_NAME }}
				(
					WordID AS CONVERT(nvarchar(200), JSON_VALUE(Content, '$.id')),
					Verse AS CONVERT(nvarchar(200), JSON_VALUE(Content, '$.verse')),
					CoreID AS CONVERT(nvarchar(200), JSON_VALUE(Content, '$.coreid')),
					Language AS CONVERT(nvarchar(200), JSON_VALUE(Content, '$.language')),
					Lemma AS CONVERT(nvarchar(200), JSON_VALUE(Content, '$.lemma')),
					Codes AS CONVERT(nvarchar(200), JSON_VALUE(Content, '$.codes')),
					Content [nvarchar](max) NOT NULL
				);			
			END`
	createWLCIndexes = `
	CREATE CLUSTERED INDEX Index{{ TABLE_NAME }}WordID ON {{ TABLE_NAME }} (WordID);
	ALTER TABLE [dbo].{{ TABLE_NAME }} ADD CONSTRAINT {{ TABLE_NAME }}ContentJson CHECK (ISJSON(Content)=1);
	`
)

type mssqlWord struct {
	ID   int64
	Data string
}

func createConnection() (*sql.DB, error) {
	cs := os.Getenv("CS")
	connection, err := sql.Open("mssql", cs)
	if err != nil {
		return nil, err
	}
	return connection, nil
}

func getPartitionSize() int {
	return 1000
}

func validateCloudConfig() error {
	cs := os.Getenv("CS")
	if len(cs) == 0 {
		return errors.New("CS is required")
	}
	return nil
}

func postPersistWLC(tableName string) error {
	db, err := createConnection()
	if err != nil {
		return err
	}
	defer db.Close()
	sql := strings.Replace(createWLCIndexes, "{{ TABLE_NAME }}", tableName, -1)
	_, err = db.Exec(sql)
	if err != nil {
		return err
	}
	return nil
}

func postPersistGNT(tableName string) error {
	db, err := createConnection()
	if err != nil {
		return err
	}
	defer db.Close()
	sql := strings.Replace(createGNTIndexes, "{{ TABLE_NAME }}", tableName, -1)
	fmt.Println(sql)
	_, err = db.Exec(sql)
	if err != nil {
		return err
	}
	return nil
}

func prepareAndPersistWlc(tableName string, bookName string, words []wlcWord) error {
	db, err := createConnection()
	if err != nil {
		return err
	}
	defer db.Close()
	sql := strings.Replace(createWLCTable, "{{ TABLE_NAME }}", tableName, -1)
	_, err = db.Exec(sql)
	if err != nil {
		return err
	}
	var prepared []mssqlWord
	for _, word := range words {
		m, _ := json.Marshal(word)
		prepared = append(prepared, mssqlWord{
			ID:   word.SequenceID,
			Data: string(m),
		})
	}
	return partitionAndPersist(db, tableName, bookName, prepared)
}

func prepareAndPersistGnt(tableName string, bookName string, words []gntWord) error {
	db, err := createConnection()
	if err != nil {
		return err
	}
	defer db.Close()
	sql := strings.Replace(createGNTTable, "{{ TABLE_NAME }}", tableName, -1)
	_, err = db.Exec(sql)
	if err != nil {
		return err
	}
	var prepared []mssqlWord
	for _, word := range words {
		m, _ := json.Marshal(word)
		prepared = append(prepared, mssqlWord{
			ID:   word.ID,
			Data: string(m),
		})
	}
	return partitionAndPersist(db, tableName, bookName, prepared)
}

func partitionAndPersist(db *sql.DB, tableName string, bookName string, prepared []mssqlWord) error {
	partitionSize := getPartitionSize()
	fmt.Printf("partition size: %d\n", partitionSize)
	segmentNumber := 1
	fmt.Printf("Saving %s (%d words)...\n", bookName, len(prepared))
	for idxRange := range partition(len(prepared), partitionSize) {
		// fmt.Printf("partition: %d %d %d\n", idxRange.Low, idxRange.High, idxRange.High-idxRange.Low)
		segment := prepared[idxRange.Low:idxRange.High]
		err := persist(db, tableName, segment)
		if err != nil {
			return err
		}
		percent := (float64(segmentNumber) * float64((partitionSize)) / float64(len(prepared))) * 100
		if percent > 100 {
			percent = 100
		}
		fmt.Printf("%s %0.2f%% complete\n", bookName, percent)
		segmentNumber++
	}
	return nil
}

func persist(db *sql.DB, tableName string, segment []mssqlWord) error {
	txn, err := db.Begin()
	if err != nil {
		return err
	}

	stmt, err := txn.Prepare(mssql.CopyIn(tableName, mssql.BulkOptions{}, "Content"))
	if err != nil {
		return err
	}

	for _, word := range segment {
		_, err = stmt.Exec(word.Data)
		if err != nil {
			return err
		}
	}

	_, err = stmt.Exec()
	if err != nil {
		return err
	}

	err = stmt.Close()
	if err != nil {
		return err
	}

	err = txn.Commit()
	if err != nil {
		return err
	}
	return nil
}
