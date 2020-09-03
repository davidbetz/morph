package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const (
	bookOffset = 60
)

type gntMorphology struct {
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

type gntWord struct {
	Verse      string        `json:"verse"`
	ID         int64         `json:"id"`
	Codes      string        `json:"codes"`
	Morphology gntMorphology `json:"morphology"`
	Text       string        `json:"text"`
	Word       string        `json:"word"`
	Normalized string        `json:"normalized"`
	Lemma      string        `json:"lemma"`
}

func (t *Gnt) getPartName(part string) string {
	switch part {
	case "A-":
		return "adjective"
	case "C-":
		return "conjunction"
	case "D-":
		return "adverb"
	case "I-":
		return "interjection"
	case "N-":
		return "noun"
	case "P-":
		return "preposition"
	case "RA":
		return "definite article"
	case "RD":
		return "demonstrative pronoun"
	case "RI":
		return "interrogative/indefinite pronoun"
	case "RP":
		return "personal pronoun"
	case "RR":
		return "relative pronoun"
	case "V-":
		return "verb"
	case "X-":
		return "particle"
	default:
		errorf("invalid part " + part)
		return ""
	}
}

func (t *Gnt) getMorphology(part string, code string) gntMorphology {
	var person string
	var tense string
	var voice string
	var mood string
	var _case string
	var number string
	var gender string
	var degree string
	for i, value := range strings.Split(code, "") {
		switch i {
		case 0:
			person = t.personLookup[value]
			break
		case 1:
			tense = t.tenseLookup[value]
			break
		case 2:
			voice = t.voiceLookup[value]
			break
		case 3:
			mood = t.moodLookup[value]
			break
		case 4:
			_case = t.caseLookup[value]
			break
		case 5:
			number = t.numberLookup[value]
			break
		case 6:
			gender = t.genderLookup[value]
			break
		case 7:
			degree = t.degreeLookup[value]
			break
		}
	}
	return gntMorphology{
		Part:   t.getPartName(part),
		Person: person,
		Tense:  tense,
		Voice:  voice,
		Mood:   mood,
		Case:   _case,
		Number: number,
		Gender: gender,
		Degree: degree,
	}
}

// Gnt represents the GNT parser
type Gnt struct {
	personLookup map[string]string
	tenseLookup  map[string]string
	voiceLookup  map[string]string
	moodLookup   map[string]string
	caseLookup   map[string]string
	numberLookup map[string]string
	genderLookup map[string]string
	degreeLookup map[string]string
	bookOrder    map[string]int
	bookNames    map[int]string
}

func (t *Gnt) setupTables() {
	t.bookOrder = make(map[string]int, 39)
	t.bookOrder["Matthew"] = 1
	t.bookOrder["Mark"] = 2
	t.bookOrder["Luke"] = 3
	t.bookOrder["John"] = 4
	t.bookOrder["Acts"] = 5
	t.bookOrder["Romans"] = 6
	t.bookOrder["1 Corinthians"] = 7
	t.bookOrder["2 Corinthians"] = 8
	t.bookOrder["Galatians"] = 9
	t.bookOrder["Ephesians"] = 10
	t.bookOrder["Philippians"] = 11
	t.bookOrder["Colossians"] = 12
	t.bookOrder["1 Thessalonians"] = 13
	t.bookOrder["2 Thessalonians"] = 14
	t.bookOrder["1 Timothy"] = 15
	t.bookOrder["2 Timothy"] = 16
	t.bookOrder["Titus"] = 17
	t.bookOrder["Philemon"] = 18
	t.bookOrder["Hebrews"] = 19
	t.bookOrder["James"] = 20
	t.bookOrder["1 Peter"] = 21
	t.bookOrder["2 Peter"] = 22
	t.bookOrder["1 John"] = 23
	t.bookOrder["2 John"] = 24
	t.bookOrder["3 John"] = 25
	t.bookOrder["Jude"] = 26
	t.bookOrder["Revelation"] = 27

	bookOffset := 60
	t.bookNames = make(map[int]string, 27)
	for k, v := range t.bookOrder {
		v += bookOffset
		t.bookNames[v] = k
	}

	t.personLookup = map[string]string{
		"1": "first",
		"2": "second",
		"3": "third",
	}
	t.tenseLookup = map[string]string{
		"P": "present",
		"I": "imperfect",
		"F": "future",
		"A": "aorist",
		"X": "perfect",
		"Y": "pluperfect",
	}
	t.voiceLookup = map[string]string{
		"A": "active",
		"M": "middle",
		"P": "passive",
	}
	t.moodLookup = map[string]string{
		"I": "indicative",
		"D": "imperative",
		"S": "subjunctive",
		"O": "optative",
		"N": "infinitive",
		"P": "participle",
	}
	t.caseLookup = map[string]string{
		"N": "nominative",
		"G": "genitive",
		"D": "dative",
		"A": "accusative",
	}
	t.numberLookup = map[string]string{
		"S": "singular",
		"P": "plural",
	}
	t.genderLookup = map[string]string{
		"M": "masculine",
		"F": "feminine",
		"N": "neuter",
	}
	t.degreeLookup = map[string]string{
		"C": "comparative",
		"S": "superlative",
	}
}

func (t *Gnt) createAbsoluteID(verse string, id int) int64 {
	newID, err := strconv.ParseInt(verse+strconv.Itoa(id), 10, 32)
	if err != nil {
		errorf(err.Error())
	}
	return newID
}

func (t *Gnt) readFile(filename string) ([][]string, error) {
	if filepath.Ext(filename) != ".txt" {
		return nil, errors.New("Skip")
	}
	debug(fmt.Sprintf("PARSING: %v\n", filename))
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	var lines [][]string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, strings.Split(scanner.Text(), " "))
	}
	return lines, nil
}

func (t *Gnt) ParseFileContent(filename string) ([]gntWord, error) {
	words := []gntWord{}
	lines, err := t.readFile(filename)
	if err != nil {
		return nil, err
	}
	var id int
	originalVerse := ""
	for _, parts := range lines {
		if originalVerse != parts[0] {
			originalVerse = parts[0]
			id = 1
		}
		bookNumber, _ := strconv.ParseInt(originalVerse[0:2], 10, 32)
		verse := strconv.Itoa(int(bookNumber)+39) + originalVerse[2:4] + originalVerse[4:6]
		uniqueID := t.createAbsoluteID(verse, id)
		words = append(words, gntWord{
			ID:         uniqueID,
			Verse:      verse,
			Codes:      parts[2],
			Morphology: t.getMorphology(parts[1], parts[2]),
			Text:       parts[3],
			Word:       parts[4],
			Normalized: parts[5],
			Lemma:      parts[6],
		})
		id++
	}
	return words, nil
}

// CreateGnt creates a Gnt parser
func CreateGnt() *Gnt {
	gnt := &Gnt{}
	gnt.setupTables()
	return gnt
}
