package parser

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/davidbetz/morph/internal/models"
	"github.com/davidbetz/morph/internal/util"
	"golang.org/x/text/unicode/norm"
)

type gntBookData struct {
	Name string
	Data []models.GntWord
}

func (t *Gnt) getBookNumber(filename string) int {
	r, _ := regexp.Compile("([0-9]+)-([a-zA-Z0-9]+)-morphgnt")
	results := r.FindStringSubmatch(filename)
	if len(results) < 2 {
		util.Errorf("Invalid filename " + filename)
	}
	bookNumber, err := strconv.ParseInt(results[1], 10, 32)
	if err != nil {
		util.Errorf(err.Error())
	}
	return int(bookNumber)
}

func (t *Gnt) readData(books chan *gntBookData) {
	pwd, _ := os.Getwd()
	location := util.LoadEnvVarDef("DATA_LOCATION", path.Join(pwd, "data"))
	folder := path.Join(location, "morphgnt")
	files, err := os.ReadDir(folder)
	if err != nil {
		util.Errorf(err.Error())
	}
	for _, f := range files {
		filename := f.Name()
		bookName := filename[0 : len(filename)-len(filepath.Ext(filename))]
		words, err := t.ParseFileContent(path.Join(folder, filename))
		if err != nil {
			if err.Error() == "Skip" {
				continue
			}
			util.Errorf(err.Error())
		}
		books <- &gntBookData{
			bookName,
			words,
		}
		fmt.Printf("Added %s\n", bookName)
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

func (t *Gnt) Process(targets models.Target) error {
	books := make(chan *gntBookData)
	go t.readData(books)
	for book := range books {
		if book == nil {
			return nil
		}
		fmt.Printf("Parsed %s. Saving...\n", book.Name)
		name := t.bookNames[t.getBookNumber(book.Name)]
		err := targets.PrepareAndPersistGnt(t.getTableName(), name, book.Data)
		if err != nil {
			return err
		}
	}
	return targets.PostPersistWLC(t.getTableName())
}

type StrongsGntMapEntry struct {
	Strongs string
	Lemma   string
	Word    string
}

type StrongsData []StrongsGntMapEntry

// func (l *StrongsData) FindWord(word string) (StrongsGntMapEntry, error) {
// 	normalizedWord := norm.NFC.String(strings.ToLower(word))
// 	// if word == "λάθρᾳ" {
// 	// 	fmt.Printf("normalizedWord: %s\n", normalizedWord)
// 	// }
// 	for _, lemma := range *l {
// 		normalizedLemma := norm.NFC.String(strings.ToLower(lemma.Word))
// 		if normalizedLemma == normalizedWord {
// 			return lemma, nil
// 		}
// 	}
// 	return StrongsGntMapEntry{}, fmt.Errorf("word not found: %s", word)
// }

func (t *Gnt) LoadStrongsMap() (StrongsData, error) {
	relativePath := "SBLGNT-add-ons/re-shape_MorphGNT/morphological-lexicon.csv"
	pwd, _ := os.Getwd()
	location := util.LoadEnvVarDef("DATA_LOCATION", path.Join(pwd, "data"))
	filename := path.Join(location, relativePath)
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	var entries []StrongsGntMapEntry
	reader := csv.NewReader(file)
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		if len(record) < 4 {
			util.Errorf("Invalid record: %v", record)
			continue
		}
		entry := StrongsGntMapEntry{
			Strongs: strings.TrimPrefix(record[0], "G"),
			Lemma:   record[1],
			Word:    record[2],
		}
		entries = append(entries, entry)
	}
	return entries, nil
}

// GreekAnalyticalEntry represents a parsed Greek lexicon entry.
type GreekAnalyticalEntry struct {
	GreekWord        string
	Gloss            string
	POS              string
	FullCitationForm string
	Forms            []Form
}

// Form represents a morphological form of the Greek word.
type Form struct {
	GreekWord string
	POS       string
	Details   string
}

func (t *Gnt) parseGreekEntry(fields []string) (*GreekAnalyticalEntry, error) {
	if len(fields) < 2 {
		return nil, errors.New("insufficient fields: expected at least 2")
	}

	entry := &GreekAnalyticalEntry{}
	text := fields[1]

	// Regex patterns
	greekWordPattern := regexp.MustCompile(`<h2><grk>([^<]+)</grk></h2>`)
	glossPattern := regexp.MustCompile(`gloss: <font color='red'>([^<]+)</font>`)
	posPattern := regexp.MustCompile(`<font color='purple'>([^<]+)</font>`)
	citationPattern := regexp.MustCompile(`full-citation-form: <font color='blue'><grk>([^<]+)</grk></font>`)
	formPattern := regexp.MustCompile(`<grk>([^<]+)</grk>｜<font color='purple'>([^<]+)</font>｜<i>([^<]+)</i>`)

	// Extract Greek word
	// fmt.Printf("Parsing text: %s\n", text)
	if matches := greekWordPattern.FindStringSubmatch(text); matches != nil {
		match := matches[1]
		match = strings.Split(match, "/")[0]
		entry.GreekWord = norm.NFC.String(strings.ToLower(match))
	} else {
		fmt.Println("----------------")
		fmt.Printf("Pattern: %s\n", greekWordPattern)
		fmt.Printf("Base text: %s\n", text)
		return nil, fmt.Errorf("no Greek word found in text: %s", text)
	}

	// Extract gloss
	if matches := glossPattern.FindStringSubmatch(text); matches != nil {
		entry.Gloss = matches[1]
	}

	// Extract POS (first occurrence only)
	if matches := posPattern.FindStringSubmatch(text); matches != nil {
		entry.POS = matches[1]
	}

	// Extract full citation form
	if matches := citationPattern.FindStringSubmatch(text); matches != nil {
		entry.FullCitationForm = matches[1]
	}

	// Extract forms
	forms := formPattern.FindAllStringSubmatch(text, -1)
	for _, match := range forms {
		entry.Forms = append(entry.Forms, Form{
			GreekWord: norm.NFC.String(strings.ToLower(match[1])),
			POS:       match[2],
			Details:   match[3],
		})
	}

	return entry, nil
}

// parseCSV reads and parses entries from a CSV file or similar input.
func (t *Gnt) parseGreekEntries() (GreekEntryData, error) {
	relativePath := "SBLGNT-add-ons/end-user-modules/e-Sword/files/dic_MorphGNT_plus_analytical.csv"
	pwd, _ := os.Getwd()
	location := util.LoadEnvVarDef("DATA_LOCATION", path.Join(pwd, "data"))
	filename := path.Join(location, relativePath)
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = '\t' // Assuming tab-delimited input
	reader.TrimLeadingSpace = true

	var entries []*GreekAnalyticalEntry
	line := 0

	for {
		line++
		fields, err := reader.Read()
		if err != nil {
			break
		}

		if len(fields) < 2 {
			fmt.Printf("Skipping line %d: wrong number of fields\n", line)
			continue
		}

		entry, err := t.parseGreekEntry(fields)
		if err != nil {
			fmt.Printf("Error parsing line %d: %v\n", line, err)
			continue
		}

		entries = append(entries, entry)
	}
	return entries, nil
}

// func (t *Gnt) parseGreekEntries() (GreekEntryData, error) {
// 	relativePath := "SBLGNT-add-ons/end-user-modules/e-Sword/files/dic_MorphGNT_plus_analytical.csv"
// 	pwd, _ := os.Getwd()
// 	location := util.LoadEnvVarDef("DATA_LOCATION", path.Join(pwd, "data"))
// 	filename := path.Join(location, relativePath)
// 	file, err := os.Open(filename)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer file.Close()

// 	var entries GreekEntryData
// 	reader := csv.NewReader(file)
// 	for {
// 		record, err := reader.Read()
// 		if err == io.EOF {
// 			break
// 		}
// 		if err != nil {
// 			return nil, err
// 		}
// 		if len(record) < 2 {
// 			continue
// 		}
// 		fmt.Printf("Column 1: %s, Column 2: %s\n", record[0], record[1])
// 		entry, err := t.parseGreekEntry(record[1])
// 		//+ use the first column instead of the HTML parsed one
// 		entry.GreekWord = record[0]
// 		if err != nil {
// 			return nil, err
// 		}
// 		entries = append(entries, entry)
// 	}
// 	return entries, nil
// }

type GreekEntryData []*GreekAnalyticalEntry

func (l *GreekEntryData) SaveToDisk() error {
	filename := path.Join("./output", "greek-analytical-entries") + ".txt"
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}
	defer file.Close()
	for _, entry := range *l {
		_, err = file.WriteString(fmt.Sprintf("%s|%s|%s|%s\n", entry.GreekWord, entry.Gloss, entry.POS, entry.FullCitationForm))
		if err != nil {
			return err
		}
	}
	return nil
}

func (l *GreekEntryData) FindWord(word models.GntWord) (*GreekAnalyticalEntry, error) {
	// if word == "λάθρᾳ" {
	// 	fmt.Printf("normalizedWord: %s\n", normalizedWord)
	// }
	for _, lemma := range *l {
		if lemma.GreekWord == word.NormalizedLemma {
			return lemma, nil
		}
	}
	// panic(fmt.Sprintf("Word not found: %s|%s", word.NormalizedLemma, word.Lemma))
	return &GreekAnalyticalEntry{}, fmt.Errorf("word not found: %s", word.Lemma)
}

func (t *Gnt) Count(targets models.Target) error {
	entries, err := t.parseGreekEntries()
	if err != nil {
		return err
	}
	// err = entries.SaveToDisk()
	// if err != nil {
	// 	return err
	// }
	// for _, entry := range entries {
	// 	fmt.Printf("Greek word: %s\n", entry.GreekWord)
	// }
	// filename := path.Join("./output", "strongs-greek") + ".gob"
	// var lemmas models.StrongsData
	// file, err := os.Open(filename)
	// if err != nil {
	// 	return err
	// }
	// defer file.Close()
	// decoder := gob.NewDecoder(file)
	// err = decoder.Decode(&lemmas)
	// if err != nil {
	// 	return err
	// }
	bucketLimits := []int{1, 10, 50, 100}
	books := make(chan *gntBookData)
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
			foundWord, err := entries.FindWord(word)
			if err != nil {
				t.AddWordToNotFoundTextFile(word.Lemma)
			}
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
	var allWordCountsMap = make(map[string]models.WordCount)
	for _, bookBuckets := range bookBuckets {
		for _, data := range bookBuckets {
			for word, wordCount := range data {
				if existingWordCount, exists := allWordCountsMap[word]; exists {
					existingWordCount.Count += wordCount.Count
					allWordCountsMap[word] = existingWordCount
				} else {
					allWordCountsMap[word] = wordCount
				}
			}
		}
	}

	var allWordCounts []models.WordCount
	for _, wordCount := range allWordCountsMap {
		allWordCounts = append(allWordCounts, wordCount)
	}
	overallBuckets := CreateBuckets(allWordCounts, bucketLimits)
	// overallBuckets := make(map[int]map[string]models.WordCount)
	// for _, bookBuckets := range bookBuckets {
	// 	for bucketKey, data := range bookBuckets {
	// 		if _, exists := overallBuckets[bucketKey]; !exists {
	// 			overallBuckets[bucketKey] = make(map[string]models.WordCount)
	// 		}
	// 		for word, wordCount := range data {
	// 			if _, exists := overallBuckets[bucketKey][word]; !exists {
	// 				overallBuckets[bucketKey][word] = wordCount
	// 			} else {
	// 				overallBuckets[bucketKey][word] = models.WordCount{
	// 					Word:  word,
	// 					Count: overallBuckets[bucketKey][word].Count + wordCount.Count,
	// 					Gloss: wordCount.Gloss,
	// 				}
	// 			}
	// 		}
	// 	}
	// }
	return targets.PersistCounts("overall", overallBuckets)
}

func CreateBuckets(wordCounts []models.WordCount, bucketLimits []int) map[int]map[string]models.WordCount {
	buckets := make(map[int]map[string]models.WordCount)
	for _, word := range wordCounts {
		var bucketKey int
		for _, limit := range bucketLimits {
			if word.Count <= limit {
				bucketKey = limit
				break
			}
		}
		if _, exists := buckets[bucketKey]; !exists {
			buckets[bucketKey] = make(map[string]models.WordCount)
		}
		buckets[bucketKey][word.Word] = word
	}
	return buckets
}

func (t *Gnt) AddWordToNotFoundTextFile(word string) {
	fmt.Printf("AddWordToNotFoundTextFile::Word not found: %s\n", word)
	folder := path.Join("./output")
	if _, err := os.Stat(folder); os.IsNotExist(err) {
		_ = os.MkdirAll(folder, 0777)
	}
	filename := path.Join(folder, "not-found") + ".txt"
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		util.Errorf(err.Error())
	}
	defer file.Close()
	_, err = file.WriteString(word + "\n")
	if err != nil {
		util.Errorf(err.Error())
	}
}

func (t *Gnt) Render(targets models.Target) error {
	books := make(chan *gntBookData)
	go t.readData(books)
	for book := range books {
		if book == nil {
			return nil
		}
		var words []string
		for _, word := range book.Data {
			words = append(words, word.Word)
		}
		err := targets.PersistRender(book.Name, words)
		if err != nil {
			return err
		}
	}
	return nil
}
