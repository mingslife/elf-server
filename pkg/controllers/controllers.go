package controllers

import (
	"elf-server/pkg/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func HandleError(ctx *gin.Context, msg string) {
	ctx.JSON(http.StatusInternalServerError, gin.H{"message": msg})
}

func NoRoute(ctx *gin.Context) {
	settings := models.GetAllSettingsMap()
	navigations := models.GetAllNavigationsActive()

	ctx.HTML(http.StatusNotFound, "error.jet", gin.H{
		"Settings":    settings,
		"Navigations": navigations,
		"Title":       settings["app.title"],
		"Keywords":    settings["app.keywords"],
		"Description": settings["app.description"],
	})
}
