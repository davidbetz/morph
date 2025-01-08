package models

type Parser interface {
	Process(targets Target) error
	Count(targets Target) error
	Render(targets Target) error
}

type Target interface {
	ValidateCloudConfig() error
	PrepareAndPersistWlc(tableName string, bookName string, words []WlcWord) error
	PrepareAndPersistGnt(tableName string, bookName string, words []GntWord) error
	PostPersistWLC(tableName string) error
	PostPersistGNT(tableName string) error
	PersistRender(bookName string, words []string) error
	PersistCounts(bookName string, counts map[int]map[string]WordCount) error
	PersistStrongs(language string, data []Lemma) error
	PersistMacula(language string, data []MaculaWord) error
}
