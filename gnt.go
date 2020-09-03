package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
)

type gntBookData struct {
	Name string
	Data []gntWord
}

func (t *Gnt) getBookNumber(filename string) int {
	r, _ := regexp.Compile("([0-9]+)-([a-zA-Z0-9]+)-morphgnt")
	results := r.FindStringSubmatch(filename)
	if len(results) < 2 {
		errorf("Invalid filename " + filename)
	}
	bookNumber, err := strconv.ParseInt(results[1], 10, 32)
	if err != nil {
		errorf(err.Error())
	}
	return int(bookNumber)
}

func (t *Gnt) readData(books chan *gntBookData) {
	folder := os.Getenv("SOURCE")
	if len(folder) == 0 {
		folder = "./morphgnt/"
	}
	files, err := ioutil.ReadDir(folder)
	if err != nil {
		errorf(err.Error())
	}
	for _, f := range files {
		filename := f.Name()
		bookName := filename[0 : len(filename)-len(filepath.Ext(filename))]
		words, err := t.ParseFileContent(path.Join(folder, filename))
		if err != nil {
			if err.Error() == "Skip" {
				continue
			}
			errorf(err.Error())
		}
		books <- &gntBookData{
			bookName,
			words,
		}
	}
	close(books)
}

func (t *Gnt) getTableName() string {
	tableName := os.Getenv("TABLE_NAME")
	if len(tableName) == 0 {
		tableName = "morphgnt"
	}
	return tableName
}

func (t *Gnt) Process() error {
	books := make(chan *gntBookData)
	go t.readData(books)
	for book := range books {
		if book == nil {
			return nil
		}
		fmt.Printf("Parsed %s. Saving...\n", book.Name)
		name := t.bookNames[t.getBookNumber(book.Name)]
		err := prepareAndPersistGnt(t.getTableName(), name, book.Data)
		if err != nil {
			return err
		}
	}
	return postPersistWLC(t.getTableName())
}
