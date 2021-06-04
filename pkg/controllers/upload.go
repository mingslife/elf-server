package controllers

import (
	"fmt"
	"image"
	"net/http"
	"os"
	"path"
	"regexp"
	"strconv"

	"github.com/disintegration/imaging"
	"github.com/gin-gonic/gin"

	"elf-server/pkg/utils"
)

type UploadController struct{}

func (c *UploadController) UploadFile(ctx *gin.Context) {
	fileName, _, _, _ := c.saveUploadedFile(ctx)

	path := c.getPath(fileName)

	ctx.JSON(http.StatusCreated, gin.H{"path": path})
}

func (c *UploadController) UploadImage(ctx *gin.Context) {
	fileName, filePath, suffix, _ := c.saveUploadedFile(ctx)
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
		fileName = utils.NewUUID() + ".png"
		imagePath = path.Join("upload", fileName)
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

	path := c.getPath(fileName)

	ctx.JSON(http.StatusCreated, gin.H{"path": path})
}

func (c *UploadController) saveUploadedFile(ctx *gin.Context) (fileName, filePath, suffix string, fileSize int64) {
	file, _ := ctx.FormFile("file")
	suffix = path.Ext(file.Filename)
	fileName = utils.NewUUID() + suffix
	filePath = path.Join("upload", fileName)
	fileSize = file.Size
	ctx.SaveUploadedFile(file, filePath)
	return
}

func (c *UploadController) getPath(fileName string) string {
	return fmt.Sprintf("/upload/%s", fileName)
}

func NewUploadController(r gin.IRouter) *UploadController {
	c := &UploadController{}
	r.Group("/upload").
		POST("/file", c.UploadFile).
		POST("/image", c.UploadImage).
		OPTIONS("/:route", nil)
	return c
}
