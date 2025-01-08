package models

import "fmt"

type GntMorphology struct {
	Part   string `json:"part,omitempty"`
	Person string `json:"person,omitempty"`
	Tense  string `json:"tense,omitempty"`
	Voice  string `json:"voice,omitempty"`
	Mood   string `json:"mood,omitempty"`
	Case   string `json:"case,omitempty"`
	Number string `json:"number,omitempty"`
	Gender string `json:"gender,omitempty"`
	Degree string `json:"degree,omitempty"`
}

type GntWord struct {
	Verse           string        `json:"verse"`
	ID              int64         `json:"id"`
	Codes           string        `json:"codes"`
	Morphology      GntMorphology `json:"morphology"`
	Text            string        `json:"text"`
	Word            string        `json:"word"`
	Normalized      string        `json:"normalized"`
	Lemma           string        `json:"lemma"`
	NormalizedLemma string        `json:"normalized_lemma"`
}

type WlcWord struct {
	Codes            string              `json:"codes"`
	Language         string              `json:"language"`
	Lemma            string              `json:"lemma"`
	ID               string              `json:"coreid"`
	Morphology       []map[string]string `json:"morphology"`
	SequenceID       int64               `json:"id"`
	Verse            string              `json:"verse"`
	MorphologyString string
}

type WordCount struct {
	ID       string `json:"id"`
	Word     string `json:"word"`
	Count    int    `json:"count"`
	Gloss    string `json:"gloss"`
	Category string `json:"category"`
}

func (t *WordCount) String() string {
	return fmt.Sprint(t.Count, ":", t.Word, ":", t.Gloss)
}

func CreateWordCountBuckets(wordCounts []WordCount, bucketLimits []int) map[int]map[string]WordCount {
	fmt.Printf("Creating word count buckets\n")
	buckets := make(map[int]map[string]WordCount)
	for _, word := range wordCounts {
		// fmt.Printf("Word: %s, Count: %d\n", word.Word, word.Count)
		var bucketKey int
		for _, limit := range bucketLimits {
			if word.Count <= limit {
				bucketKey = limit
				break
			}
		}
		if _, exists := buckets[bucketKey]; !exists {
			buckets[bucketKey] = make(map[string]WordCount)
		}
		buckets[bucketKey][word.Word] = word
	}
	return buckets
}
