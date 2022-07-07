package main

import (
	"github.com/mingslife/bone"

	"elf-server/pkg/component"
	"elf-server/pkg/conf"
	"elf-server/pkg/module/api"
	"elf-server/pkg/module/auth"
	"elf-server/pkg/module/category"
	"elf-server/pkg/module/comment"
	"elf-server/pkg/module/log"
	"elf-server/pkg/module/navigation"
	"elf-server/pkg/module/portal"
	"elf-server/pkg/module/post"
	"elf-server/pkg/module/poststatistics"
	"elf-server/pkg/module/reader"
	"elf-server/pkg/module/setting"
	"elf-server/pkg/module/upload"
	"elf-server/pkg/module/user"
)

var cfg *conf.Config

func main() {
	app := bone.NewApplication(&bone.ApplicationOptions{
		Host: cfg.Host,
		Port: cfg.Port,
	})
	app.Use(
		new(component.Tracker),
		new(component.Logger),
		new(component.Database),
		new(component.Jet),
		new(component.Cors),
		new(component.Jwt),
		new(component.Cache),
		new(component.Captcha),
		new(component.Log),
		new(upload.Module),
		new(setting.Module),
		new(user.Module),
		new(category.Module),
		new(navigation.Module),
		new(reader.Module),
		new(poststatistics.Module),
		new(post.Module),
		new(comment.Module),
		new(auth.Module),
		new(log.Module),
		new(api.Module),
		new(portal.Module),
	)

	app.Run()
}

func init() {
	cfg = conf.ParserConfigFromEnv()
}
