package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"unicode/utf8"
)

type node struct {
	Lookup    map[string]string
	Name      string
	TopName   string
	Next      *node
	Alternate *node
	Decider   func(string) bool
}

type wlcWord struct {
	Codes            string              `json:"codes"`
	Language         string              `json:"language"`
	Lemma            string              `json:"lemma"`
	ID               string              `json:"coreid"`
	Morphology       []map[string]string `json:"morphology"`
	SequenceID       int64               `json:"id"`
	Verse            string              `json:"verse"`
	MorphologyString string
}

func shiftString(text string) string {
	_, i := utf8.DecodeRuneInString(text)
	text = text[i:]
	return text
}

type Wlc struct {
	bookOrder                  map[string]int
	filenames                  map[string]string
	partOfSpeechLookup         map[string]string
	hebrewStemLookup           map[string]string
	aramaicVerbLookup          map[string]string
	verbConjugationTypesLookup map[string]string
	adjectiveLookup            map[string]string
	nounLookup                 map[string]string
	pronounLookup              map[string]string
	prepositionLookup          map[string]string
	suffixLookup               map[string]string
	particleLookup             map[string]string
	hebrewPersonLookup         map[string]string
	hebrewGenderLookup         map[string]string
	hebrewNumberLookup         map[string]string
	stateLookup                map[string]string
	languageLookup             map[string]string
	notUsedLookup              map[string]string
	languageVerbLookup         map[string]map[string]string
	trees                      map[string]*node
}

func (t *Wlc) setupTables() {
	t.partOfSpeechLookup = make(map[string]string, 9)
	t.hebrewStemLookup = make(map[string]string, 27)
	t.aramaicVerbLookup = make(map[string]string, 26)
	t.verbConjugationTypesLookup = make(map[string]string, 11)
	t.adjectiveLookup = make(map[string]string, 4)
	t.nounLookup = make(map[string]string, 3)
	t.pronounLookup = make(map[string]string, 5)
	t.prepositionLookup = make(map[string]string, 1)
	t.suffixLookup = make(map[string]string, 4)
	t.particleLookup = make(map[string]string, 9)
	t.hebrewPersonLookup = make(map[string]string, 3)
	t.hebrewGenderLookup = make(map[string]string, 4)
	t.hebrewNumberLookup = make(map[string]string, 3)
	t.stateLookup = make(map[string]string, 3)
	t.languageLookup = make(map[string]string, 2)
	t.languageVerbLookup = make(map[string]map[string]string, 2)
	t.notUsedLookup = make(map[string]string, 1)

	t.bookOrder = make(map[string]int, 39)
	t.filenames = make(map[string]string, 39)

	t.bookOrder["Genesis"] = 1
	t.bookOrder["Exodus"] = 2
	t.bookOrder["Leviticus"] = 3
	t.bookOrder["Numbers"] = 4
	t.bookOrder["Deuteronomy"] = 5
	t.bookOrder["Joshua"] = 6
	t.bookOrder["Judges"] = 7
	t.bookOrder["Ruth"] = 8
	t.bookOrder["I Samuel"] = 9
	t.bookOrder["II Samuel"] = 10
	t.bookOrder["I Kings"] = 11
	t.bookOrder["II Kings"] = 12
	t.bookOrder["I Chronicles"] = 13
	t.bookOrder["II Chronicles"] = 14
	t.bookOrder["Ezra"] = 15
	t.bookOrder["Nehemiah"] = 16
	t.bookOrder["Esther"] = 17
	t.bookOrder["Job"] = 18
	t.bookOrder["Psalms"] = 19
	t.bookOrder["Proverbs"] = 20
	t.bookOrder["Ecclesiastes"] = 21
	t.bookOrder["Song of Solomon"] = 22
	t.bookOrder["Isaiah"] = 23
	t.bookOrder["Jeremiah"] = 24
	t.bookOrder["Lamentations"] = 25
	t.bookOrder["Ezekiel"] = 26
	t.bookOrder["Daniel"] = 27
	t.bookOrder["Hosea"] = 28
	t.bookOrder["Joel"] = 29
	t.bookOrder["Amos"] = 30
	t.bookOrder["Obadiah"] = 31
	t.bookOrder["Jonah"] = 32
	t.bookOrder["Micah"] = 33
	t.bookOrder["Nahum"] = 34
	t.bookOrder["Habakkuk"] = 35
	t.bookOrder["Zephaniah"] = 36
	t.bookOrder["Haggai"] = 37
	t.bookOrder["Zechariah"] = 38
	t.bookOrder["Malachi"] = 39

	t.partOfSpeechLookup["A"] = "adjective"
	t.partOfSpeechLookup["C"] = "conjunction"
	t.partOfSpeechLookup["D"] = "adverb"
	t.partOfSpeechLookup["N"] = "noun"
	t.partOfSpeechLookup["P"] = "pronoun"
	t.partOfSpeechLookup["R"] = "preposition"
	t.partOfSpeechLookup["S"] = "suffix"
	t.partOfSpeechLookup["T"] = "particle"
	t.partOfSpeechLookup["V"] = "verb"

	t.hebrewStemLookup["q"] = "qal"
	t.hebrewStemLookup["N"] = "niphal"
	t.hebrewStemLookup["p"] = "piel"
	t.hebrewStemLookup["P"] = "pual"
	t.hebrewStemLookup["h"] = "hiphil"
	t.hebrewStemLookup["H"] = "hophal"
	t.hebrewStemLookup["t"] = "hithpael"
	t.hebrewStemLookup["o"] = "polel"
	t.hebrewStemLookup["O"] = "polal"
	t.hebrewStemLookup["r"] = "hithpolel"
	t.hebrewStemLookup["m"] = "poel"
	t.hebrewStemLookup["M"] = "poal"
	t.hebrewStemLookup["k"] = "palel"
	t.hebrewStemLookup["K"] = "pulal"
	t.hebrewStemLookup["Q"] = "qal passive"
	t.hebrewStemLookup["l"] = "pilpel"
	t.hebrewStemLookup["L"] = "polpal"
	t.hebrewStemLookup["f"] = "hithpalpel"
	t.hebrewStemLookup["D"] = "nithpael"
	t.hebrewStemLookup["j"] = "pealal"
	t.hebrewStemLookup["i"] = "pilel"
	t.hebrewStemLookup["u"] = "hothpaal"
	t.hebrewStemLookup["c"] = "tiphil"
	t.hebrewStemLookup["v"] = "hishtaphel"
	t.hebrewStemLookup["w"] = "nithpalel"
	t.hebrewStemLookup["y"] = "nithpoel"
	t.hebrewStemLookup["z"] = "hithpoel"

	t.aramaicVerbLookup["q"] = "peal"
	t.aramaicVerbLookup["Q"] = "peil"
	t.aramaicVerbLookup["u"] = "hithpeel"
	t.aramaicVerbLookup["p"] = "pael"
	t.aramaicVerbLookup["P"] = "ithpaal"
	t.aramaicVerbLookup["M"] = "hithpaal"
	t.aramaicVerbLookup["a"] = "aphel"
	t.aramaicVerbLookup["h"] = "haphel"
	t.aramaicVerbLookup["s"] = "saphel"
	t.aramaicVerbLookup["e"] = "shaphel"
	t.aramaicVerbLookup["H"] = "hophal"
	t.aramaicVerbLookup["i"] = "ithpeel"
	t.aramaicVerbLookup["t"] = "hishtaphel"
	t.aramaicVerbLookup["v"] = "ishtaphel"
	t.aramaicVerbLookup["w"] = "hithaphel"
	t.aramaicVerbLookup["o"] = "polel"
	t.aramaicVerbLookup["z"] = "ithpoel"
	t.aramaicVerbLookup["r"] = "hithpolel"
	t.aramaicVerbLookup["f"] = "hithpalpel"
	t.aramaicVerbLookup["b"] = "hephal"
	t.aramaicVerbLookup["c"] = "tiphel"
	t.aramaicVerbLookup["m"] = "poel"
	t.aramaicVerbLookup["l"] = "palpel"
	t.aramaicVerbLookup["L"] = "ithpalpel"
	t.aramaicVerbLookup["O"] = "ithpolel"
	t.aramaicVerbLookup["G"] = "ittaphal"

	t.verbConjugationTypesLookup["p"] = "perfect (qatal)"
	t.verbConjugationTypesLookup["q"] = "sequential perfect (weqatal)"
	t.verbConjugationTypesLookup["i"] = "imperfect (yiqtol)"
	t.verbConjugationTypesLookup["w"] = "sequential imperfect (wayyiqtol)"
	t.verbConjugationTypesLookup["h"] = "cohortative"
	t.verbConjugationTypesLookup["j"] = "jussive"
	t.verbConjugationTypesLookup["v"] = "imperative"
	t.verbConjugationTypesLookup["r"] = "participle active"
	t.verbConjugationTypesLookup["s"] = "participle passive"
	t.verbConjugationTypesLookup["a"] = "infinitive absolute"
	t.verbConjugationTypesLookup["c"] = "infinitive construct"

	t.adjectiveLookup["a"] = "adjective"
	t.adjectiveLookup["c"] = "cardinal number"
	t.adjectiveLookup["g"] = "gentilic"
	t.adjectiveLookup["o"] = "ordinal number"

	t.nounLookup["c"] = "common"
	t.nounLookup["g"] = "gentilic"
	t.nounLookup["p"] = "proper name"

	t.pronounLookup["d"] = "demonstrative"
	t.pronounLookup["f"] = "indefinite"
	t.pronounLookup["i"] = "interrogative"
	t.pronounLookup["p"] = "personal"
	t.pronounLookup["r"] = "relative"

	t.prepositionLookup["d"] = "definite article"

	t.suffixLookup["d"] = "directional he"
	t.suffixLookup["h"] = "paragogic he"
	t.suffixLookup["n"] = "paragogic nun"
	t.suffixLookup["p"] = "pronominal"

	t.particleLookup["a"] = "affirmation"
	t.particleLookup["d"] = "definite article"
	t.particleLookup["e"] = "exhortation"
	t.particleLookup["i"] = "interrogative"
	t.particleLookup["j"] = "interjection"
	t.particleLookup["m"] = "demonstrative"
	t.particleLookup["n"] = "negative"
	t.particleLookup["o"] = "direct object marker"
	t.particleLookup["r"] = "relative"

	t.hebrewPersonLookup["1"] = "first"
	t.hebrewPersonLookup["2"] = "second"
	t.hebrewPersonLookup["3"] = "third"

	t.hebrewGenderLookup["b"] = "both (noun)"
	t.hebrewGenderLookup["c"] = "common (verb)"
	t.hebrewGenderLookup["f"] = "feminine"
	t.hebrewGenderLookup["m"] = "masculine"

	t.hebrewNumberLookup["d"] = "dual"
	t.hebrewNumberLookup["p"] = "plural"
	t.hebrewNumberLookup["s"] = "singular"

	t.stateLookup["a"] = "absolute"
	t.stateLookup["c"] = "construct"
	t.stateLookup["d"] = "determined"

	t.languageLookup["H"] = "Hebrew"
	t.languageLookup["A"] = "Aramaic"

	t.languageVerbLookup["H"] = t.hebrewStemLookup
	t.languageVerbLookup["A"] = t.aramaicVerbLookup

	t.notUsedLookup["x"] = "-"
	for k := range t.bookOrder {
		t.filenames[strings.ToLower(strings.Replace(k, " ", "", -1))] = k
	}
	//+ POS to *node
	t.trees = make(map[string]*node, 7)
	gnsBranch := &node{
		Lookup: t.hebrewGenderLookup,
		Name:   "Gender",
		Next: &node{
			Lookup: t.hebrewNumberLookup,
			Name:   "Number",
			Next: &node{
				Lookup: t.stateLookup,
				Name:   "State",
			},
		},
	}
	pgnBranch := &node{
		Lookup: t.hebrewPersonLookup,
		Name:   "Person",
		Next: &node{
			Lookup: t.hebrewGenderLookup,
			Name:   "Gender",
			Next: &node{
				Lookup: t.hebrewNumberLookup,
				Name:   "Number",
			},
		},
	}
	t.trees["C"] = &node{
		Name: "conjunction",
	}
	t.trees["D"] = &node{
		Name: "adverb",
	}
	t.trees["A"] = &node{
		Lookup:  t.adjectiveLookup,
		TopName: "adjective",
		Name:    "Type",
		Next:    gnsBranch,
	}
	t.trees["N"] = &node{
		Lookup:  t.nounLookup,
		TopName: "noun",
		Name:    "Type",
		Next:    gnsBranch,
	}
	t.trees["P"] = &node{
		Lookup:  t.pronounLookup,
		TopName: "pronoun",
		Name:    "Type",
		Decider: func(next string) bool {
			return next == "x"
		},
		Next: pgnBranch,
		Alternate: &node{
			Lookup: t.notUsedLookup,
			Name:   "-",
			Next: &node{
				Lookup: t.hebrewGenderLookup,
				Name:   "Gender",
				Next: &node{
					Lookup: t.hebrewNumberLookup,
					Name:   "Number",
					Next: &node{
						Lookup: t.stateLookup,
						Name:   "State",
					},
				},
			},
		},
	}
	t.trees["R"] = &node{
		Name:   "preposition",
		Lookup: t.prepositionLookup,
	}
	t.trees["S"] = &node{
		Lookup:  t.suffixLookup,
		TopName: "suffix",
		Name:    "Type",
		Next:    pgnBranch,
	}
	t.trees["T"] = &node{
		TopName: "particle",
		Name:    "Type",
		Lookup:  t.particleLookup,
	}
	verbDecider := func(next string) bool {
		return next != "r" && next != "s"
	}
	t.trees["VH"] = &node{
		Lookup:  t.hebrewStemLookup,
		TopName: "verb",
		Name:    "Stem",
		Next: &node{
			Lookup:  t.verbConjugationTypesLookup,
			Name:    "Conjugation",
			Decider: verbDecider,
			Next: &node{
				Lookup: t.hebrewPersonLookup,
				Name:   "Person",
				Next:   gnsBranch,
			},
			Alternate: &node{
				Lookup: t.hebrewGenderLookup,
				Name:   "Gender",
				Next: &node{
					Lookup: t.hebrewNumberLookup,
					Name:   "Number",
					Next: &node{
						Lookup: t.stateLookup,
						Name:   "State",
					},
				},
			},
		},
	}
	t.trees["VA"] = &node{
		Lookup:  t.aramaicVerbLookup,
		TopName: "Verb",
		Name:    "Stem",
		Next: &node{
			Lookup:  t.verbConjugationTypesLookup,
			Name:    "Conjugation",
			Decider: verbDecider,
			Next: &node{
				Lookup: t.hebrewPersonLookup,
				Name:   "Person",
				Next:   gnsBranch,
			},
			Alternate: &node{
				Lookup: t.hebrewGenderLookup,
				Name:   "Gender",
				Next: &node{
					Lookup: t.hebrewNumberLookup,
					Name:   "Number",
					Next: &node{
						Lookup: t.stateLookup,
						Name:   "State",
					},
				},
			},
		},
	}
}

func (t *Wlc) parseMorphology(morph string) (string, []map[string]string) {
	original := morph
	languageCode := string(morph[0])
	morph = shiftString(morph)
	var language string
	switch languageCode {
	case "H":
		language = "Hebrew"
		break
	case "A":
		language = "Aramaic"
		break
	}
	debug(fmt.Sprintf("STARTING NEXT WORD %s %s\n", original, language))
	var morphologyArray []map[string]string
	parts := strings.Split(morph, "/")
	for _, part := range parts {
		// original := part
		m := make(map[string]string)
		partOfSpeechCode := string(part[0])
		part = shiftString(part)
		var tree *node
		if partOfSpeechCode == "V" {
			tree = t.trees[partOfSpeechCode+languageCode]
			debug(fmt.Sprintf("TREE: %s %v\n", partOfSpeechCode+languageCode, tree))
		} else {
			tree = t.trees[partOfSpeechCode]
			debug(fmt.Sprintf("TREE: %s %v\n", partOfSpeechCode, tree))
		}
		topName := tree.TopName
		if len(topName) == 0 {
			topName = tree.Name
		}
		m["Part"] = topName
		debug(fmt.Sprintf("\tSTARTING NEXT PART, %s %s\n", original, tree.Name))
		for _, l := range part {
			letter := string(l)
			debug(fmt.Sprintf("\t\tSTARTING NEXT LETTER, %s %q\n", letter, tree))
			if tree.Name != "-" {
				m[tree.Name] = tree.Lookup[letter]
			}
			if tree != nil {
				if d := tree.Decider; d != nil {
					if d(letter) {
						debug(fmt.Sprintf("DECIDER CALLED TRUE %s\n", letter))
						tree = tree.Next
					} else {
						debug(fmt.Sprintf("DECIDER CALLED FALSE %s\n", letter))
						tree = tree.Alternate
					}
				} else {
					tree = tree.Next
				}
			}
		}
		morphologyArray = append(morphologyArray, m)
	}
	return language, morphologyArray
}

func (t *Wlc) Parse(word []string, verseID string, sequence int64) wlcWord {
	lemma := word[0]
	id := word[1]
	morph := word[2]
	language, morphologyArray := t.parseMorphology(morph)
	var outer []string
	for _, morph := range morphologyArray {
		var inner []string
		for k, v := range morph {
			inner = append(inner, k+"="+v)
		}
		outer = append(outer, strings.Join(inner, ","))
	}
	return wlcWord{
		Codes:            morph,
		Language:         language,
		Morphology:       morphologyArray,
		MorphologyString: strings.Join(outer, "|"),
		Lemma:            lemma,
		ID:               id,
		Verse:            verseID,
		SequenceID:       sequence,
	}
}

func (t *Wlc) ReadFile(filename string) ([]byte, error) {
	if filepath.Ext(filename) != ".json" {
		return nil, errors.New("not json")
	}
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return ioutil.ReadAll(file)
}

func (t *Wlc) ParseFileContent(bookName string, filename string) ([]wlcWord, error) {
	buffer, err := t.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	var obj [][][][]string
	err = json.Unmarshal(buffer, &obj)
	if err != nil {
		return nil, err
	}
	bookID := t.bookOrder[bookName]
	wordid := 1
	var words []wlcWord
	for ci, chapter := range obj {
		for vi, verse := range chapter {
			for _, word := range verse {
				verseID := fmt.Sprintf("%02d%03d%03d", bookID, ci+1, vi+1)
				sequenceID := fmt.Sprintf("%02d%03d%03d%03d", bookID, ci+1, vi+1, wordid)
				sequence, err := strconv.ParseInt(sequenceID, 10, 64)
				if err != nil {
					panic(err)
				}
				word := t.Parse(word, verseID, sequence)
				words = append(words, word)
				wordid++
			}
		}
	}
	return words, nil
}

func CreateWlc() *Wlc {
	wlc := &Wlc{}
	wlc.setupTables()
	return wlc
}
