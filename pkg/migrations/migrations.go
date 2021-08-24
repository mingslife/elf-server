package migrations

import (
	"log"

	"elf-server/pkg/models"
)

type migrationFunc func()

var migrationFuncs map[uint64]migrationFunc = map[uint64]migrationFunc{}

func ExecuteMigrations() {
	nextVersion := uint64(1)
	if migration := models.GetMigrationOfMaxVersion(); migration != nil {
		nextVersion = migration.Version + 1
	}

	for {
		if f, ok := migrationFuncs[nextVersion]; ok {
			f()
			(&models.Migration{Version: nextVersion}).Save()
			log.Printf("Migration version %d executed.", nextVersion)
			nextVersion++
		} else {
			break
		}
	}
}
