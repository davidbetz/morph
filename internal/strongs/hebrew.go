package strongs

import (
	"encoding/xml"
	"fmt"
	"log"
	"os"
	"path"
	"sort"
	"strconv"
	"sync" // Add sync package

	"github.com/davidbetz/morph/internal/models"
	"github.com/davidbetz/morph/internal/util"
)

type Hebrew struct {
	lemmas map[int]models.Lemma
	mu     sync.Mutex // Add a mutex to the struct
}

func (t *Hebrew) saveEntries(entryChan <-chan models.Lemma) {
	for entry := range entryChan {
		t.mu.Lock() // Lock the mutex before writing to the map
		if _, exists := t.lemmas[entry.Number]; !exists {
			t.lemmas[entry.Number] = entry
			fmt.Printf("Processed Entry: %s\n", entry.String())
		}
		t.mu.Unlock() // Unlock the mutex after writing to the map
	}
}

func (t *Hebrew) Count(targets models.Target) error {
	return nil
}

func (t *Hebrew) Process(targets models.Target) error {
	entryChan := make(chan models.Lemma)
	t.lemmas = make(map[int]models.Lemma)
	go t.saveEntries(entryChan)
	pwd, _ := os.Getwd()
	location := util.LoadEnvVarDef("DATA_LOCATION", path.Join(pwd, "data"))
	filePath := path.Join(location, "strongs/hebrew/StrongHebrewG.xml")
	if err := t.processEntries(filePath, entryChan); err != nil {
		fmt.Printf("Error: %v\n", err)
	}
	var lemmas []models.Lemma
	t.mu.Lock() // Lock the mutex before reading from the map
	for _, lemma := range t.lemmas {
		lemmas = append(lemmas, lemma)
	}
	t.mu.Unlock() // Unlock the mutex after reading from the map
	sort.Slice(lemmas, func(i, j int) bool {
		return lemmas[i].Number < lemmas[j].Number
	})
	return targets.PersistStrongs("hebrew", lemmas)
}

func CreateHebrewParser() models.Parser {
	return &Hebrew{}
}

func Load() (*models.Osis, error) {
	xmlFile, err := os.Open("strongs_snippet.xml")
	if err != nil {
		return nil, err
	}
	defer func() {
		if cerr := xmlFile.Close(); cerr != nil {
			log.Printf("Error closing file: %v", cerr)
		}
	}()
	var osis models.Osis
	decoder := xml.NewDecoder(xmlFile)
	if err := decoder.Decode(&osis); err != nil {
		return nil, err
	}
	return &osis, nil
}

func transformEntryToLemma(entry models.Entry) models.Lemma {
	n, err := strconv.Atoi(entry.N)
	if err != nil {
		log.Printf("Error converting %s to int: %v", entry.N, err)
	}
	return models.Lemma{
		Word:   entry.Word.Lemma,
		Gloss:  entry.List.Concat(),
		Number: n,
	}
}

func GetLemmas(osis *models.Osis) map[int]models.Lemma {
	lemmas := make(map[int]models.Lemma, len(osis.OsisText.Glossary.Entries))
	for _, entry := range osis.OsisText.Glossary.Entries {
		n, err := strconv.Atoi(entry.N)
		if err != nil {
			log.Printf("Error converting %s to int: %v", entry.N, err)
			continue
		}
		firstGloss := entry.List.Items[0]
		if firstGloss == "" {
			log.Printf("No gloss found for %s", entry.Word.Value)
		}
		lemmas[n] = models.Lemma{
			Word:   entry.Word.Lemma,
			Gloss:  firstGloss,
			Number: n,
		}
	}
	return lemmas
}

func (t *Hebrew) processEntries(fileName string, entryChan chan<- models.Lemma) error {
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
	var inGlossary bool

	for {
		tok, err := decoder.Token()
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			return fmt.Errorf("error decoding XML: %w", err)
		}

		switch elem := tok.(type) {
		case xml.StartElement:
			if elem.Name.Local == "div" {
				for _, attr := range elem.Attr {
					if attr.Name.Local == "type" && attr.Value == "glossary" {
						inGlossary = true
					} else if attr.Name.Local == "type" && attr.Value == "entry" && inGlossary {
						var entry models.Entry
						if err := decoder.DecodeElement(&entry, &elem); err != nil {
							return fmt.Errorf("error decoding entry: %w", err)
						}
						entryChan <- transformEntryToLemma(entry)
					}
				}
			}
		case xml.EndElement:
			if elem.Name.Local == "div" && inGlossary {
				inGlossary = false
			}
		}
	}
	return nil
}

func (t *Hebrew) Render(targets models.Target) error {
	return nil
}
