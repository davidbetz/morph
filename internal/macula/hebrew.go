package macula

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"path"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/davidbetz/morph/internal/models"
	"github.com/davidbetz/morph/internal/util"
)

type Hebrew struct {
	words  []models.MaculaWord
	counts map[string]models.WordCount
	mu     sync.Mutex // Add a mutex to the struct
}

func (t *Hebrew) parseTSV(r io.Reader) error {
	fmt.Printf("Parsing Hebrew TSV\n")
	scanner := bufio.NewScanner(r)
	var isHeader = true

	for scanner.Scan() {
		line := scanner.Text()
		if isHeader {
			// Skip header row
			isHeader = false
			continue
		}

		// Parse the line as tab-separated values
		reader := csv.NewReader(strings.NewReader(line))
		reader.Comma = '\t'
		record, err := reader.Read()
		if err != nil {
			return err
		}

		word := models.MaculaWord{
			ID:               record[0],
			Ref:              record[1],
			Class:            record[2],
			Text:             record[3],
			Transliteration:  record[4],
			After:            record[5],
			StrongNumberX:    record[6],
			StrongLemma:      record[7],
			SenseNumber:      record[8],
			Greek:            record[9],
			GreekStrong:      record[10],
			Gloss:            record[11],
			English:          record[12],
			Mandarin:         record[13],
			Stem:             record[14],
			Morph:            record[15],
			Lang:             record[16],
			Lemma:            record[17],
			Pos:              record[18],
			Person:           record[19],
			Gender:           record[20],
			Number:           record[21],
			State:            record[22],
			Type:             record[23],
			LexDomain:        record[24],
			ContextualDomain: record[25],
			CoreDomain:       record[26],
			SDBH:             record[27],
			Extends:          record[28],
			Frame:            record[29],
			SubjRef:          record[30],
			ParticipantRef:   record[31],
		}
		t.mu.Lock()
		t.words = append(t.words, word)
		t.mu.Unlock()
	}

	if err := scanner.Err(); err != nil {
		return err
	}
	return nil
}

func CreateHebrewParser() models.Parser {
	return &Hebrew{
		counts: make(map[string]models.WordCount),
	}
}

func (t *Hebrew) Load() error {
	pwd, _ := os.Getwd()
	location := util.LoadEnvVarDef("DATA_LOCATION", path.Join(pwd, "data"))
	filename := path.Join(location, "macula-hebrew/WLC/tsv/macula-hebrew.tsv")
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("error opening file: %w", err)
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	return t.parseTSV(reader)
}

func (t *Hebrew) Process(targets models.Target) error {
	err := t.Load()
	if err != nil {
		return err
	}
	return targets.PersistMacula("hebrew", t.words)
}

func (t *Hebrew) Count(targets models.Target) error {
	start := time.Now() // Start measuring time

	if err := t.Load(); err != nil {
		return err
	}

	// Channel to process words for counting
	countCh := make(chan models.MaculaWord, 1000)

	// Populate countCh in a goroutine
	go func() {
		for _, word := range t.words {
			countCh <- word
		}
		close(countCh)
	}()

	i := 0
	fmt.Printf("Counting instances of %d words...\n", len(t.words))
	for word := range countCh {
		t.mu.Lock()
		if !strings.ContainsAny(word.StrongNumberX, "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ") {
			if wordCount, ok := t.counts[word.Lemma]; !ok {
				meaning := word.English
				if meaning == "" {
					if word.Gloss != "" {
						fmt.Printf("Meaning for %s is empty, using gloss: %s\n", word.Lemma, word.Gloss)
					}
					meaning = word.Gloss
				}
				if meaning != "" {
					t.counts[word.Lemma] = models.WordCount{
						ID:       word.ID,
						Word:     word.Lemma,
						Count:    1,
						Gloss:    meaning,
						Category: word.Class,
					}
				} else {
					fmt.Printf("Meaning for %s is empty, skipping...\n", word.Lemma)
				}
			} else {
				wordCount.Count++
				t.counts[word.Lemma] = wordCount
			}
			if i%1000 == 0 {
				fmt.Printf("Processed %d words\n", i)
			}
			i++
		}
		t.mu.Unlock()
	}
	fmt.Printf("Sorting %d words...\n", len(t.counts))
	var sortedWords []string
	for word := range t.counts {
		sortedWords = append(sortedWords, word)
	}
	sort.Slice(sortedWords, func(i, j int) bool {
		return t.counts[sortedWords[i]].Count > t.counts[sortedWords[j]].Count
	})
	bucketLimits := []int{1, 9, 49, 99}
	var words []models.WordCount
	for _, word := range sortedWords {
		words = append(words, t.counts[word])
	}
	buckets := models.CreateWordCountBuckets(words, bucketLimits)

	elapsed := time.Since(start) // Measure elapsed time
	fmt.Printf("Count function took %s\n", elapsed)

	return targets.PersistCounts("hebrew", buckets)
}

func (t *Hebrew) Render(targets models.Target) error {
	return models.NewNotImplementedError("Render")
}
