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
		panic("Invalid filename " + filename)
	}
	bookNumber, err := strconv.ParseInt(results[1], 10, 32)
	if err != nil {
		panic(err)
	}
	return int(bookNumber)
}

func (t *Gnt) readData(c chan *gntBookData, e chan error) {
	folder := sourceFileLocation
	if len(folder) == 0 {
		folder = "./morphgnt/"
	}
	files, err := ioutil.ReadDir(folder)
	if err != nil {
		e <- err
	}
	for _, f := range files {
		filename := f.Name()
		bookName := filename[0 : len(filename)-len(filepath.Ext(filename))]
		words, err := t.ParseFileContent(path.Join(folder, filename))
		if err != nil {
			if err.Error() == "Skip" {
				continue
			}
			e <- err
		}
		c <- &gntBookData{
			bookName,
			words,
		}
	}
	c <- nil
}

func (t *Gnt) getTableName() string {
	tableName := os.Getenv("TABLE_NAME")
	if len(tableName) == 0 {
		tableName = "morphgnt"
	}
	return tableName
}

func (t *Gnt) Process() error {
	c := make(chan *gntBookData)
	e := make(chan error)
	go t.readData(c, e)
	for {
		select {
		case contents := <-c:
			if contents == nil {
				return nil
			}
			fmt.Printf("Parsed %s. Saving...\n", contents.Name)
			name := t.bookNames[t.getBookNumber(contents.Name)]
			err := prepareAndPersistGnt(t.getTableName(), name, contents.Data)
			if err != nil {
				return err
			}
		case err := <-e:
			if err != nil {
				return err
			}
			return postPersistWLC(t.getTableName())
		}
	}
}
