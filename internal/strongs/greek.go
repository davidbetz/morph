package strongs

import (
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"sort"
	"strconv"
	"strings"

	"github.com/davidbetz/morph/internal/models"
	"github.com/davidbetz/morph/internal/util"
)

type Greek struct {
	lemmas map[int]models.Lemma
}

func (t *Greek) saveEntries(entryChan <-chan models.Lemma) {
	for entry := range entryChan {
		if _, exists := t.lemmas[entry.Number]; !exists {
			t.lemmas[entry.Number] = entry
			fmt.Printf("Processed Entry: %s\n", entry.String())
		}
	}
}

func (t *Greek) Count(targets models.Target) error {
	return nil
}

func (t *Greek) Process(targets models.Target) error {
	entryChan := make(chan models.Lemma)
	t.lemmas = make(map[int]models.Lemma)
	go t.saveEntries(entryChan)
	pwd, _ := os.Getwd()
	location := util.LoadEnvVarDef("DATA_LOCATION", path.Join(pwd, "data"))
	filePath := path.Join(location, "strongs/greek/StrongsGreekDictionaryXML_1.4/strongsgreek.xml")
	if err := t.processEntries(filePath, entryChan); err != nil {
		fmt.Printf("Error: %v\n", err)
	}
	var lemmas []models.Lemma
	for _, lemma := range t.lemmas {
		lemmas = append(lemmas, lemma)
	}
	sort.Slice(lemmas, func(i, j int) bool {
		return lemmas[i].Number < lemmas[j].Number
	})
	return targets.PersistStrongs("greek", lemmas)
}

func CreateGreekParser() models.Parser {
	return &Greek{}
}

func (t *Greek) processEntries(fileName string, entryChan chan<- models.Lemma) error {
	defer close(entryChan)
	xmlFile, err := os.Open(fileName)
	if err != nil {
		return fmt.Errorf("error opening file: %w", err)
	}
	defer func() {
		if cerr := xmlFile.Close(); cerr != nil {
			log.Printf("Error closing file: %v", cerr)
		}
	}()
	decoder := xml.NewDecoder(xmlFile)

	var inEntry bool
	var lemma models.Lemma
	for {
		tok, err := decoder.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("error reading XML: %w", err)
		}

		switch elem := tok.(type) {
		case xml.StartElement:
			if elem.Name.Local == "entry" {
				inEntry = true
				lemma = models.Lemma{}
				for _, attr := range elem.Attr {
					if attr.Name.Local == "strongs" {
						lemma.Number, _ = strconv.Atoi(attr.Value)
					}
				}
			}
			if inEntry {
				switch elem.Name.Local {
				case "greek":
					for _, attr := range elem.Attr {
						if attr.Name.Local == "unicode" {
							lemma.Word = attr.Value
						}
					}
				case "strongs_def":
					var strongsDef string
					if err := decoder.DecodeElement(&strongsDef, &elem); err == nil {
						lemma.Gloss = strings.Trim(strings.ReplaceAll(strongsDef, "\n", ""), " ")
					}
				}
			}
		case xml.EndElement:
			if elem.Name.Local == "entry" {
				inEntry = false
				entryChan <- lemma
			}
		}
	}
	return nil
}

func (t *Greek) Render(targets models.Target) error {
	return nil
}
