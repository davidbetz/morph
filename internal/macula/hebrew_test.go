package macula

import (
	"strings"
	"testing"

	"github.com/davidbetz/morph/internal/models"
)

func TestParseTSV(t *testing.T) {
	data := `xml:id	ref	class	text	transliteration	after	strongnumberx	stronglemma	sensenumber	greek	greekstrong	gloss	english	mandarin	stem	morph	lang	lemma	pos	person	gender	number	state	type	lexdomain	contextualdomain	coredomain	sdbh	extends	frame	subjref	participantref
o010010010011	GEN 1:1!1	prep	בְּ	bə		0871a	בְּ		ἐν	1722		in			R	H	בְּ	preposition													
o010010010012	GEN 1:1!1	noun	רֵאשִׁ֖ית	rēʾšiyṯ	 	7225	רֵאשִׁית	1	ἀρξῇ	746		beginning	起初		Ncfsa	H	רֵאשִׁית	noun		feminine	singular	absolute	common	002003003004		168	006653001001000				
o010010010021	GEN 1:1!2	verb	בָּרָ֣א	bārāʾ	 	1254	בָּרָא	1	ἐποίησεν	4160	he.created	created	创造	qal	Vqp3ms	H	בָּרָא	verb	third	masculine	singular		qatal	002002002005		028 055	001156001002000		A0:010010010031; A1:010010010052;010010010072;		`

	hebrew := &Hebrew{}
	if err := hebrew.parseTSV(strings.NewReader(data)); err != nil {
		t.Fatalf("ParseTSV failed: %v", err)
	}

	// Collect and verify results
	expected := []models.MaculaWord{
		{
			ID: "o010010010011", Ref: "GEN 1:1!1", Class: "prep", Text: "בְּ",
			Transliteration: "bə", StrongNumberX: "0871a", StrongLemma: "בְּ",
			Greek: "ἐν", GreekStrong: "1722", Gloss: "in", English: "",
		},
		{
			ID: "o010010010012", Ref: "GEN 1:1!1", Class: "noun", Text: "רֵאשִׁ֖ית",
			Transliteration: "rēʾšiyṯ", StrongNumberX: "7225", StrongLemma: "רֵאשִׁית",
			SenseNumber: "1", Greek: "ἀρξῇ", GreekStrong: "746", Gloss: "beginning",
			English: "起初", Pos: "noun", Gender: "feminine", Number: "singular",
		},
		{
			ID: "o010010010021", Ref: "GEN 1:1!2", Class: "verb", Text: "בָּרָ֣א",
			Transliteration: "bārāʾ", StrongNumberX: "1254", StrongLemma: "בָּרָא",
			SenseNumber: "1", Greek: "ἐποίησεν", GreekStrong: "4160", Gloss: "he.created",
			English: "created", Pos: "verb", Person: "third", Gender: "masculine",
			Number: "singular", State: "qatal",
		},
	}

	results := hebrew.words

	// Check the number of results
	if len(results) != len(expected) {
		t.Fatalf("Expected %d words, got %d", len(expected), len(results))
	}

	// Check each parsed word
	for i, word := range results {
		if word.ID != expected[i].ID || word.Ref != expected[i].Ref || word.Text != expected[i].Text {
			t.Errorf("Word mismatch at index %d: got %+v, expected %+v", i, word, expected[i])
		}
	}
}
