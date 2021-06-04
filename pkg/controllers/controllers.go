package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func HandleError(ctx *gin.Context, msg string) {
	ctx.JSON(http.StatusInternalServerError, gin.H{"message": msg})
}
