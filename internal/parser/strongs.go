package parser

import (
	"github.com/davidbetz/morph/internal/models"
	"github.com/davidbetz/morph/internal/strongs"
)

func CreateStrongHebrew() models.Parser {
	return strongs.CreateHebrewParser()
}

func CreateStrongGreek() models.Parser {
	return strongs.CreateGreekParser()
}
