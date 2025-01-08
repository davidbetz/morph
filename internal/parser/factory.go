package parser

import (
	"github.com/davidbetz/morph/internal/macula"
	"github.com/davidbetz/morph/internal/models"
	"github.com/davidbetz/morph/internal/strongs"
)

func Create(mode, variant string) models.Parser {
	switch mode {
	case "gnt":
		return CreateGnt()
	case "wlc":
		return CreateWlc(variant)
	case "strongs-hebrew":
		return strongs.CreateHebrewParser()
	case "strongs-greek":
		return strongs.CreateGreekParser()
	case "macula-hebrew":
		return macula.CreateHebrewParser()
	default:
		return nil
	}
}
