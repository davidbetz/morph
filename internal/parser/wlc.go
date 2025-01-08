package parser

import (
	"encoding/gob"
	"fmt"
	"os"
	"path"
	"sort"
	"strings"

	"github.com/davidbetz/morph/internal/models"
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
	pwd, _ := os.Getwd()
	location := util.LoadEnvVarDef("DATA_LOCATION", path.Join(pwd, "data"))
	folder := path.Join(location, "morphwlc")
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

func (t *Wlc) Process(targets models.Target) error {
	books := make(chan *wlcBookData)
	go t.readData(books)
	for book := range books {
		fmt.Printf("Parsed %s. Saving...\n", book.Name)
		err := targets.PrepareAndPersistWlc(t.getTableName(), book.Name, book.Data)
		if err != nil {
			return err
		}
	}
	return targets.PostPersistWLC(t.getTableName())
}

func stripHebrewVowels(input string) string {
	vowels := []rune{'ְ', 'ֱ', 'ֲ', 'ֳ', 'ִ', 'ֵ', 'ֶ', 'ַ', 'ָ', 'ֹ', 'ֺ', 'ֻ', 'ּ', 'ֽ', 'ׁ', 'ׂ', 'ׄ'}
	output := []rune{}
	for _, r := range input {
		isVowel := false
		for _, v := range vowels {
			if r == v {
				isVowel = true
				break
			}
		}
		if !isVowel {
			output = append(output, r)
		}
	}
	return string(output)
}

var strip = []string{"ה", "ו", "ב", "ל"}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func (t *Wlc) Count(targets models.Target) error {
	filename := path.Join("./output", "strongs-hebrew") + ".gob"
	var lemmas models.StrongsData
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("error opening file: %s. If gob dos not exist run 'morph -mode strongs-hebrew -target gob'", filename)
	}
	defer file.Close()
	decoder := gob.NewDecoder(file)
	err = decoder.Decode(&lemmas)
	if err != nil {
		return err
	}
	books := make(chan *wlcBookData)
	go t.readData(books)
	bookBuckets := make(map[string]map[int]map[string]models.WordCount)
	for book := range books {
		if book == nil {
			return nil
		}
		wordCounts := make(map[string]models.WordCount)
		fmt.Printf("Counting %s...\n", book.Name)
		fmt.Printf("Total words: %d\n", len(book.Data))
		for _, word := range book.Data {
			// word.Lemma = stripHebrewVowels(word.Lemma)
			// if idx := strings.LastIndex(word.Lemma, "/"); idx != -1 {
			// if len(word.Lemma) > idx+1 && contains(strip, string(word.Lemma[idx+1])) {
			lemma := word.Lemma
			foundWord, err := lemmas.FindWord(word.Lemma)
			if err == nil {
				panic(fmt.Sprintf("Word not found: %s", word.Lemma))
			}
			// word.Lemma = word.Lemma[:idx]
			fmt.Printf("CLEANED Word: %s => %s\n", lemma, word.Lemma)
			// 	}
			// }
			if wordCount, ok := wordCounts[word.Lemma]; !ok {
				wordCounts[word.Lemma] = models.WordCount{
					Word:  word.Lemma,
					Count: 1,
					Gloss: foundWord.Gloss,
				}
			} else {
				wordCount.Count++
				wordCounts[word.Lemma] = wordCount
			}
		}
		var sortedWords []string
		for word := range wordCounts {
			sortedWords = append(sortedWords, word)
		}
		sort.Slice(sortedWords, func(i, j int) bool {
			return wordCounts[sortedWords[i]].Count > wordCounts[sortedWords[j]].Count
		})
		// buckets := make(map[int]map[string]models.WordCount)
		bucketLimits := []int{1, 10, 50, 100}
		// for _, word := range sortedWords {
		// 	count := wordCounts[word]
		// 	var bucketKey int
		// 	for _, limit := range bucketLimits {
		// 		if count <= limit {
		// 			bucketKey = limit
		// 			break
		// 		}
		// 	}
		// 	if _, exists := buckets[bucketKey]; !exists {
		// 		buckets[bucketKey] = make(map[string]models.WordCount)
		// 	}
		// 	buckets[bucketKey][word] = count
		// }
		var words []models.WordCount
		for _, word := range sortedWords {
			words = append(words, wordCounts[word])
		}
		buckets := CreateBuckets(words, bucketLimits)
		bookBuckets[book.Name] = buckets
		err := targets.PersistCounts(book.Name, buckets)
		if err != nil {
			return err
		}
	}
	// overallBuckets := make(map[int]map[string]int)
	// for _, bookBuckets := range bookBuckets {
	// 	for bucketKey, data := range bookBuckets {
	// 		if _, exists := overallBuckets[bucketKey]; !exists {
	// 			overallBuckets[bucketKey] = make(map[string]int)
	// 		}
	// 		for word, count := range data {
	// 			overallBuckets[bucketKey][word] += count
	// 		}
	// 	}
	// }
	// return targets.PersistCounts("overall", overallBuckets)
	return nil
}

func (t *Wlc) Render(targets models.Target) error {
	books := make(chan *wlcBookData)
	go t.readData(books)
	for book := range books {
		if book == nil {
			return nil
		}
		var words []string
		for _, word := range book.Data {
			words = append(words, word.Lemma)
		}
		err := targets.PersistRender(book.Name, words)
		if err != nil {
			return err
		}
	}
	return nil
}
