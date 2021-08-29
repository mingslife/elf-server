package controllers

import (
	"elf-server/pkg/conf"
	"elf-server/pkg/middleware"
	"elf-server/pkg/utils"
	"net/http"
	"path/filepath"
	"regexp"
	"strconv"

	"github.com/gin-gonic/gin"
)

type LogController struct{}

func (c *LogController) GetDates(ctx *gin.Context) {
	if !conf.GetConfig().Log {
		ctx.JSON(http.StatusBadRequest, gin.H{})
		return
	}

	month := ctx.Query("month")
	if ok, err := regexp.Match("^[0-9]{4}-[0-9]{2}$", []byte(month)); !ok || err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{})
		return
	}

	var dates []string
	path := middleware.LOG.Path()
	for day := 1; day <= 31; day++ {
		dayStr := strconv.Itoa(day)
		if day < 10 {
			dayStr = "0" + dayStr
		}

		ymd := month + "-" + dayStr
		fileName := middleware.LOG.LogFileName(ymd)
		filePath := filepath.Join(path, fileName)
		if utils.IsFileExists(filePath) {
			dates = append(dates, ymd)
		}
	}
	ctx.JSON(http.StatusOK, gin.H{"dates": dates})
}

func (c *LogController) GetRawLog(ctx *gin.Context) {
	if !conf.GetConfig().Log {
		ctx.JSON(http.StatusBadRequest, gin.H{})
		return
	}

	date := ctx.Query("date")
	if ok, err := regexp.Match("^[0-9]{4}-[0-9]{2}-[0-9]{2}$", []byte(date)); !ok || err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{})
		return
	}

	path := middleware.LOG.Path()
	ymd := date
	fileName := middleware.LOG.LogFileName(ymd)
	filePath := filepath.Join(path, fileName)
	ctx.File(filePath)
}

func NewLogController(r gin.IRouter) *LogController {
	c := &LogController{}
	r.Group("/log").
		GET("/dates", c.GetDates).
		GET("/raw", c.GetRawLog).
		OPTIONS("/:route", nil)
	return c
}
