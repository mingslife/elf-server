package migrations

import "elf-server/pkg/models"

func version1() {
	models.DB.AutoMigrate(&models.Migration{})
	models.DB.AutoMigrate(&models.Setting{})
	models.DB.AutoMigrate(&models.Navigation{})
	models.DB.AutoMigrate(&models.User{})
	models.DB.AutoMigrate(&models.Category{})
	models.DB.AutoMigrate(&models.Post{})
	models.DB.AutoMigrate(&models.Comment{})
}

func init() {
	migrationFuncs[1] = version1
}
