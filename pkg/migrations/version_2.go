package migrations

import (
	"log"
	"time"

	"elf-server/pkg/models"
	"elf-server/pkg/utils"
)

func version2() {
	now := time.Now()

	(&models.Setting{
		SettingKey:   "app.title",
		SettingValue: "ELF",
		SettingTag:   "Title",
		IsPublic:     false,
	}).Save()
	(&models.Setting{
		SettingKey:   "app.keywords",
		SettingValue: "blogging,go,elf",
		SettingTag:   "Keywords",
		IsPublic:     false,
	}).Save()
	(&models.Setting{
		SettingKey:   "app.description",
		SettingValue: "A blogging system written in Go",
		SettingTag:   "Description",
		IsPublic:     false,
	}).Save()
	(&models.Setting{
		SettingKey:   "app.language",
		SettingValue: "en",
		SettingTag:   "Language",
		IsPublic:     true,
	}).Save()
	(&models.Setting{
		SettingKey:   "app.footer",
		SettingValue: "&copy; 2021 ELFIN",
		SettingTag:   "Footer",
		IsPublic:     false,
	}).Save()
	(&models.Setting{
		SettingKey:   "app.limit",
		SettingValue: "10",
		SettingTag:   "Limit",
		IsPublic:     false,
	}).Save()
	(&models.Setting{
		SettingKey:   "app.comment",
		SettingValue: "true",
		SettingTag:   "Comment",
		IsPublic:     false,
	}).Save()
	(&models.Setting{
		SettingKey:   "app.script",
		SettingValue: "console.log('Powered by ELF')",
		SettingTag:   "Script",
		IsPublic:     false,
	}).Save()
	(&models.Setting{
		SettingKey:   "app.register",
		SettingValue: "true",
		SettingTag:   "Register",
		IsPublic:     true,
	}).Save()
	(&models.Setting{
		SettingKey:   "app.inviteCode",
		SettingValue: utils.RandString(16),
		SettingTag:   "Invite Code",
		IsPublic:     false,
	}).Save()
	(&models.Setting{
		SettingKey:   "app.autoSave",
		SettingValue: "0",
		SettingTag:   "Auto Save Duration",
		IsPublic:     true,
	}).Save()

	adminPassword := utils.RandString(8)
	(&models.User{
		Username:     "admin",
		Password:     adminPassword,
		Nickname:     "Admin",
		Email:        "admin@elf",
		Phone:        "0000",
		Tags:         "admin",
		Introduction: "Adminstrator account",
		IsActive:     true,
		ActiveAt:     &now,
		Avatar:       "/assets/avatar.svg",
		Gender:       models.UserGenderMale,
		Birthday:     &now,
		Role:         models.UserRoleAdmin,
	}).Save()
	log.Printf("admin password: %s\n", adminPassword)
}

func init() {
	migrationFuncs[2] = version2
}
