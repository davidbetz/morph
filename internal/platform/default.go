//go:build !json && !aws && !azure && !gcp && !mssql && !print
// +build !json,!aws,!azure,!gcp,!mssql,!print

package platform

import (
	"errors"

	"github.com/davidbetz/morph/internal/models"
)

func ValidateCloudConfig() error {
	return errors.New("no cloud configuration specified")
}

func PrepareAndPersistWlc(tableName string, bookName string, words []models.WlcWord) error {
	return errors.New("no cloud configuration specified")
}

func PrepareAndPersistGnt(tableName string, bookName string, words []models.GntWord) error {
	return errors.New("no cloud configuration specified")
}

func PostPersistWLC(tableName string) error {
	return nil
}

func PostPersistGNT(tableName string) error {
	return nil
}
