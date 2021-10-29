package migrations

import (
	"elf-server/pkg/models"
)

func version5() {
	models.DB.AutoMigrate(&models.Reader{})
	(&models.Setting{
		SettingKey:   "app.reader",
		SettingValue: "false",
		SettingTag:   "Reader Support",
		IsPublic:     false,
	}).Save()
	(&models.Setting{
		SettingKey:   "app.smtp",
		SettingValue: "true",
		SettingTag:   "SMTP Support",
		IsPublic:     false,
	}).Save()
	(&models.Setting{
		SettingKey:   "app.smtpHost",
		SettingValue: "smtp.example.com",
		SettingTag:   "SMTP Server Host",
		IsPublic:     false,
	}).Save()
	(&models.Setting{
		SettingKey:   "app.smtpPort",
		SettingValue: "25",
		SettingTag:   "SMTP Server Port",
		IsPublic:     false,
	}).Save()
	(&models.Setting{
		SettingKey:   "app.smtpEmail",
		SettingValue: "noreply@example.com",
		SettingTag:   "SMTP User Email",
		IsPublic:     false,
	}).Save()
	(&models.Setting{
		SettingKey:   "app.smtpPassword",
		SettingValue: "",
		SettingTag:   "SMTP User Password",
		IsPublic:     false,
	}).Save()
	settingAppIcon := models.GetSettingByKey("app.icon")
	settingAppIcon.SettingValue = "/assets/elf.ico"
	settingAppIcon.Update()
	settingAppleIcon := models.GetSettingByKey("app.icon")
	settingAppleIcon.SettingValue = "/assets/elf.png"
	settingAppleIcon.Update()
	models.DB.Exec("alter table comments rename to comments_old")
	models.DB.AutoMigrate(&models.Comment{})
}

func init() {
	migrationFuncs[5] = version5
}
