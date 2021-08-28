package migrations

import (
	"elf-server/pkg/models"
	"elf-server/pkg/utils"
)

func version3() {
	models.DB.AutoMigrate(&models.Post{})
	models.DB.AutoMigrate(&models.PostStatistics{})

	var posts []*models.Post
	models.DB.Model(&models.Post{}).
		Select([]string{"id"}).
		Where("unique_id is null").
		Find(&posts)

	for _, post := range posts {
		uniqueID := utils.NewUUID()
		models.DB.Model(&models.Post{ID: post.ID}).Updates(map[string]interface{}{
			"unique_id": uniqueID,
		})
		models.DB.Save(&models.PostStatistics{UniqueID: uniqueID})
	}
}

func init() {
	migrationFuncs[3] = version3
}
