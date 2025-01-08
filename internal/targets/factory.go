package targets

import "github.com/davidbetz/morph/internal/models"

func Create(name string) models.Target {
	switch name {
	case "jsonl":
		return CreateJsonl()
	case "text":
		return CreateText()
	case "aws":
		return CreateAws()
	case "azure":
		return CreateAzure()
	case "gcp":
		return CreateGcp()
	case "print":
		return CreatePrint()
	case "mssql":
		return CreateMsSql()
	case "gob":
		return CreateGob()
	}
	return CreateText()
}
