package models

import (
	"encoding/xml"
	"fmt"
	"regexp"
	"strings"

	"golang.org/x/text/unicode/norm"
)

type Osis struct {
	XMLName  xml.Name `xml:"osis"`
	OsisText OsisText `xml:"osisText"`
}

type OsisText struct {
	OsisIDWork string      `xml:"osisIDWork,attr"`
	Lang       string      `xml:"xml:lang,attr"`
	Header     Header      `xml:"header"`
	Glossary   GlossaryDiv `xml:"div"`
}

type Header struct {
	RevisionDescs []RevisionDesc `xml:"revisionDesc"`
	Works         []Work         `xml:"work"`
	WorkPrefixes  []WorkPrefix   `xml:"workPrefix"`
}

type RevisionDesc struct {
	Resp string `xml:"resp,attr"`
	Date string `xml:"date"`
	P    string `xml:"p"`
}

type Work struct {
	OsisWork     string     `xml:"osisWork,attr"`
	Lang         string     `xml:"xml:lang,attr"`
	Title        []Title    `xml:"title"`
	Contributors []string   `xml:"contributor"`
	Creators     []Creator  `xml:"creator"`
	Dates        []Date     `xml:"date"`
	Description  string     `xml:"description"`
	Publisher    string     `xml:"publisher"`
	Identifier   Identifier `xml:"identifier"`
	Source       string     `xml:"source"`
	Rights       []Rights   `xml:"rights"`
}

type Title struct {
	Type  string `xml:"type,attr,omitempty"`
	Value string `xml:",chardata"`
}

type Creator struct {
	Role string `xml:"role,attr"`
	Name string `xml:",chardata"`
}

type Date struct {
	Event string `xml:"event,attr"`
	Type  string `xml:"type,attr,omitempty"`
	Value string `xml:",chardata"`
}

type Identifier struct {
	Type  string `xml:"type,attr"`
	Value string `xml:",chardata"`
}

type Rights struct {
	Type  string `xml:"type,attr,omitempty"`
	Value string `xml:",chardata"`
}

type WorkPrefix struct {
	Path     string `xml:"path,attr"`
	OsisWork string `xml:"osisWork,attr"`
}

type GlossaryDiv struct {
	Entries []Entry `xml:"div"`
}

type Entry struct {
	Type    string    `xml:"type,attr"`
	N       string    `xml:"n,attr"`
	Word    Word      `xml:"w"`
	Foreign []Foreign `xml:"foreign"`
	List    List      `xml:"list"`
	Notes   []Note    `xml:"note"`
}

type Word struct {
	Gloss string `xml:"gloss,attr"`
	Lemma string `xml:"lemma,attr"`
	Morph string `xml:"morph,attr"`
	POS   string `xml:"POS,attr"`
	Xlit  string `xml:"xlit,attr"`
	ID    string `xml:"ID,attr"`
	Lang  string `xml:"xml:lang,attr"`
	Value string `xml:",chardata"`
}

type Foreign struct {
	Lang  string `xml:"xml:lang,attr"`
	Words []Word `xml:"w"`
}

type List struct {
	Items []string `xml:"item"`
}

func stripBeforeParen(input string) string {
	// Define the regular expression
	re := regexp.MustCompile(`(?s)^.*?\)\s`)
	// Replace everything up to the first ") " (including it)
	result := re.ReplaceAllString(input, "")
	return result
}

func (e *List) Concat() string {
	var items = make([]string, len(e.Items))
	for i, item := range e.Items {
		items[i] = strings.Replace(stripBeforeParen(item), "  ", " ", -1)
	}
	return strings.Join(items, "; ")
}

type Note struct {
	Type  string `xml:"type,attr"`
	Value string `xml:",chardata"`
}

type Lemma struct {
	Word   string
	Gloss  string
	Number int
}

type StrongsData []Lemma

func (l *StrongsData) FindWord(word string) (Lemma, error) {
	normalizedWord := norm.NFC.String(strings.ToLower(word))
	// if word == "λάθρᾳ" {
	// 	fmt.Printf("normalizedWord: %s\n", normalizedWord)
	// }
	for _, lemma := range *l {
		normalizedLemma := norm.NFC.String(strings.ToLower(lemma.Word))
		if normalizedLemma == normalizedWord {
			return lemma, nil
		}
	}
	return Lemma{}, fmt.Errorf("word not found: %s", word)
}

func (l *Lemma) String() string {
	return fmt.Sprintf("%s: %s", l.Word, l.Gloss)
}
