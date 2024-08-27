package models

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
	Verse      string        `json:"verse"`
	ID         int64         `json:"id"`
	Codes      string        `json:"codes"`
	Morphology GntMorphology `json:"morphology"`
	Text       string        `json:"text"`
	Word       string        `json:"word"`
	Normalized string        `json:"normalized"`
	Lemma      string        `json:"lemma"`
}
