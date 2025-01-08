package strongs

import (
	"testing"
)

func TestOsisParsing(t *testing.T) {
	osis, err := Load()
	if err != nil {
		t.Fatalf("Error loading OSIS: %v", err)
	}

	t.Logf("Parsed OSIS: %+v", osis)
}

func TestEnsureContent(t *testing.T) {
	osis, err := Load()
	if err != nil {
		t.Fatalf("Error loading OSIS: %v", err)
	}

	if len(osis.OsisText.Glossary.Entries) != 2 {
		t.Fatalf("Expected 2 GlossEntries, got %d", len(osis.OsisText.Glossary.Entries))
	}
	if osis.OsisText.Glossary.Entries[0].Word.Lemma != "אָב" {
		t.Fatalf("Expected Gloss 'אָב', got %s", osis.OsisText.Glossary.Entries[0].Word.Gloss)
	}
	t.Logf("Parsed OSIS: %+v", osis)
}

func TestGetLemmas(t *testing.T) {
	osis, err := Load()
	if err != nil {
		t.Fatalf("Error loading OSIS: %v", err)
	}
	lemmas := GetLemmas(osis)
	if len(lemmas) != 2 {
		t.Fatalf("Expected 2 lemmas, got %d", len(lemmas))
	}
	if lemmas[1].Word != "אָב" {
		t.Fatalf("Expected Word 'אָב', got %s", lemmas[1].Word)
	}
}
