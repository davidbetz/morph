package main

import (
	"fmt"
	"os"
	"path"
	"strings"
)

type wlcBookData struct {
	Name string
	Data []wlcWord
}

func (t *Wlc) getTableName() string {
	tableName := os.Getenv("TABLE_NAME")
	if len(tableName) == 0 {
		tableName = "morphwlc"
	}
	return tableName
}

func (t *Wlc) getStyle() string {
	style := os.Getenv("VERSE_MODE")
	if style == "english" {
		style = "remapped"
		fmt.Println("Using English verses.")
	} else {
		style = "hebrew"
		fmt.Println("Using Hebrew verses. Specify VERSE_MODE=english for the other mode.")
	}
	return style
}

func (t *Wlc) readData(c chan *wlcBookData, e chan error) {
	folder := sourceFileLocation
	style := t.getStyle()
	if len(folder) == 0 {
		folder = fmt.Sprintf("./morphhb/%s/", style)
	}
	for n := 1; n < 40; n++ {
		var bookName string
		for name, number := range t.bookOrder {
			if number == n {
				bookName = name
			}
		}
		filename := strings.ToLower(strings.Replace(bookName, " ", "", -1))
		words, err := t.ParseFileContent(bookName, path.Join(folder, filename+".json"))
		if err != nil {
			if err.Error() == "Skip" {
				continue
			}
			panic(err)
			e <- err
		}
		c <- &wlcBookData{
			bookName,
			words,
		}
	}
	c <- nil
}

func (t *Wlc) Process() error {
	c := make(chan *wlcBookData)
	e := make(chan error)
	go t.readData(c, e)
	for {
		select {
		case contents := <-c:
			if contents == nil {
				return nil
			}
			fmt.Printf("Parsed %s. Saving...\n", contents.Name)
			err := prepareAndPersistWlc(t.getTableName(), contents.Name, contents.Data)
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
