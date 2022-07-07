package service

import (
	"image"
	"io/fs"
	"io/ioutil"
	"os"
	"path"
	"time"

	"github.com/disintegration/imaging"
	"github.com/mingslife/bone"

	"elf-server/pkg/utils"
)

type UploadService struct{}

func (s *UploadService) SaveFile(name string, data []byte) (filePath string, err error) {
	_, _, filePath, err = s.saveFile(name, data)
	if err != nil {
		return "", err
	}
	return
}

func (s *UploadService) SaveImage(name string, data []byte, compress bool, width, height int) (filePath string, err error) {
	var fileDir string
	fileDir, _, filePath, err = s.saveFile(name, data)
	if err != nil {
		return "", err
	}

	imagePath := filePath
	if compress {
		fileName := utils.NewUUID() + ".png"
		imagePath = path.Join(fileDir, fileName)

		var src image.Image
		src, err = imaging.Open(filePath, imaging.AutoOrientation(true))
		if err != nil {
			return
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
	return imagePath, nil
}

func (s *UploadService) saveFile(name string, data []byte) (fileDir, fileName, filePath string, err error) {
	fileDir = s.getUploadFileDir()
	fileName = s.generateNewFileName(name)
	filePath = path.Join(fileDir, fileName)
	err = ioutil.WriteFile(filePath, data, fs.ModePerm)
	return
}

func (*UploadService) generateNewFileName(fileName string) string {
	return utils.NewUUID() + path.Ext(fileName)
}

func (*UploadService) getUploadFileDir() string {
	fileDir := path.Join("upload", time.Now().Local().Format("20060102"))
	if !utils.IsFileExists(fileDir) {
		os.MkdirAll(fileDir, os.ModePerm)
	}
	return fileDir
}

var _ bone.Service = (*UploadService)(nil)
