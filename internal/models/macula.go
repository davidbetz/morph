package models

type MaculaWord struct {
	ID               string
	Ref              string `json:"ref"`
	Class            string `json:"class"`
	Text             string `json:"text"`
	Transliteration  string
	After            string
	StrongNumberX    string `json:"strongNumberX"`
	StrongLemma      string
	SenseNumber      string
	Greek            string
	GreekStrong      string
	Gloss            string `json:"gloss"`
	English          string `json:"english"`
	Mandarin         string
	Stem             string
	Morph            string
	Lang             string
	Lemma            string `json:"lemma"`
	Pos              string
	Person           string
	Gender           string
	Number           string
	State            string
	Type             string
	LexDomain        string
	ContextualDomain string
	CoreDomain       string
	SDBH             string
	Extends          string
	Frame            string
	SubjRef          string
	ParticipantRef   string
}
