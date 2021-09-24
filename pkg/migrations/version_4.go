package migrations

import (
	"elf-server/pkg/models"
)

func version4() {
	(&models.Setting{
		SettingKey:   "app.brand",
		SettingValue: "ELF",
		SettingTag:   "Brand",
		IsPublic:     false,
	}).Save()
	(&models.Setting{
		SettingKey:   "app.icon",
		SettingValue: "assets/elf.ico",
		SettingTag:   "Icon",
		IsPublic:     false,
	}).Save()
	(&models.Setting{
		SettingKey:   "app.appleIcon",
		SettingValue: "assets/elf.png",
		SettingTag:   "Icon for Apple",
		IsPublic:     false,
	}).Save()
}

func init() {
	migrationFuncs[4] = version4
}
