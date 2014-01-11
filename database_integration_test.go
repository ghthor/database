package database

import (
	"github.com/ghthor/database/config"
	"github.com/ghthor/gospec"
	"log"
)

var cfg config.Config

func init() {
	var err error
	cfg, err = config.ReadFromFile("config.json")
	if err != nil {
		log.Fatalf("Error reading config: %v", err)
	}
}

func DescribeDatabaseIntegration(c gospec.Context) {
}
