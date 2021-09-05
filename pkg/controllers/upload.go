package controllers

import (
	"image"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/disintegration/imaging"
	"github.com/gin-gonic/gin"

	"elf-server/pkg/utils"
)

type UploadController struct{}

func (c *UploadController) UploadFile(ctx *gin.Context) {
	_, filePath, _, _ := c.saveUploadedFile(ctx)

	path := c.getPath(filePath)

	ctx.JSON(http.StatusCreated, gin.H{"path": path})
}

func (c *UploadController) UploadImage(ctx *gin.Context) {
	_, filePath, suffix, _ := c.saveUploadedFile(ctx)
	if matched, _ := regexp.MatchString("jpg|jpeg|png|gif|tif|tiff|bmp", suffix); !matched {
		os.Remove(filePath)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Unsupport image format",
		})
		return
	}

	imagePath := filePath
	compress, _ := strconv.ParseBool(ctx.DefaultPostForm("compress", "true"))
	width, _ := strconv.Atoi(ctx.DefaultPostForm("width", "1000"))
	height, _ := strconv.Atoi(ctx.DefaultPostForm("height", "1000"))

	if compress {
		fileName := utils.NewUUID() + ".png"
		fileDir := c.getUploadFileDir()
		imagePath = path.Join(fileDir, fileName)
		src, err := imaging.Open(filePath, imaging.AutoOrientation(true))
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{})
		}
		b := src.Bounds()
		imageWidth, imageHeight := b.Max.X, b.Max.Y
		var dst *image.NRGBA
		if imageWidth > width || imageHeight > height {
			dst = imaging.Fit(src, width, height, imaging.Lanczos)
		} else {
			dst = imaging.Fit(src, imageWidth, imageHeight, imaging.Lanczos)
		}
		imaging.Save(dst, imagePath) // png.DefaultCompression
		os.Remove(filePath)
	}

	path := c.getPath(imagePath)

	ctx.JSON(http.StatusCreated, gin.H{"path": path})
}

func (c *UploadController) saveUploadedFile(ctx *gin.Context) (fileName, filePath, suffix string, fileSize int64) {
	file, _ := ctx.FormFile("file")
	suffix = path.Ext(file.Filename)
	fileName = utils.NewUUID() + suffix
	fileDir := c.getUploadFileDir()
	filePath = path.Join(fileDir, fileName)
	fileSize = file.Size
	ctx.SaveUploadedFile(file, filePath)
	return
}

func (c *UploadController) getPath(filePath string) string {
	path := strings.ReplaceAll(filepath.Clean("/"+filePath), `\`, `/`)
	return path
}

func (c *UploadController) getUploadFileDir() string {
	fileDir := path.Join("upload", time.Now().Local().Format("20060102"))
	if !utils.IsFileExists(fileDir) {
		os.MkdirAll(fileDir, os.ModePerm)
	}
	return fileDir
}

func NewUploadController(r gin.IRouter) *UploadController {
	c := &UploadController{}
	r.Group("/upload").
		POST("/file", c.UploadFile).
		POST("/image", c.UploadImage).
		OPTIONS("/:route", nil)
	return c
}
