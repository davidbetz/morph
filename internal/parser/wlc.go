package parser

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/davidbetz/morph/internal/models"
	"github.com/davidbetz/morph/internal/platform"
	"github.com/davidbetz/morph/internal/util"
)

type wlcBookData struct {
	Name string
	Data []models.WlcWord
}

func (t *Wlc) getTableName() string {
	tableName := os.Getenv("TABLE_NAME")
	if len(tableName) == 0 {
		tableName = "morphwlc"
	}
	return tableName
}

func (t *Wlc) cleanStyle() string {
	style := t.style
	if style == "english" {
		style = "remapped"
		fmt.Println("Using English verses.")
	} else {
		style = "hebrew"
		fmt.Println("Using Hebrew verses. Specify -style=english for the other mode.")
	}
	return style
}

func (t *Wlc) readData(books chan *wlcBookData) {
	folder := os.Getenv("SOURCE")
	if len(folder) == 0 {
		folder = "./morphhb/"
	}
	style := t.cleanStyle()
	folder = path.Join(folder, style)
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
			util.Errorf(err.Error())
		}
		books <- &wlcBookData{
			bookName,
			words,
		}
	}
	close(books)
}

func (t *Wlc) Process() error {
	books := make(chan *wlcBookData)
	go t.readData(books)
	for book := range books {
		fmt.Printf("Parsed %s. Saving...\n", book.Name)
		err := platform.PrepareAndPersistWlc(t.getTableName(), book.Name, book.Data)
		if err != nil {
			return err
		}
	}
	return platform.PostPersistWLC(t.getTableName())
}
