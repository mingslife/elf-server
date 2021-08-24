package main

import (
	"fmt"
	"net/http"

	"elf-server/pkg/conf"
	"elf-server/pkg/migrations"
	"elf-server/pkg/models"
	"elf-server/pkg/router"
)

func main() {
	cfg := conf.ParserConfigFromEnv()

	models.InitDB(cfg.DbHost, cfg.DbPort, cfg.DbUser, cfg.DbPwd, cfg.DbName, cfg.Debug)
	migrations.ExecuteMigrations()

	router := router.NewRouter(cfg)

	http.ListenAndServe(fmt.Sprintf("%s:%d", cfg.Host, cfg.Port), router)
}
