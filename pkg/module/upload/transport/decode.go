package transport

import (
	"context"
	"errors"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"path"
	"regexp"

	"elf-server/pkg/module/upload/model"

	"github.com/gogf/gf/util/gconv"
)

type UploadDecoder struct{}

func (d *UploadDecoder) UploadFile(_ context.Context, r *http.Request) (any, error) {
	f, fh, err := r.FormFile("file")
	if err != nil {
		return nil, err
	}
	defer f.Close()

	fileName, fileExt, fileSize, e := d.parseFileHeader(fh)
	if e != nil {
		return nil, e
	}
	if fileSize > 50*1024*1024 { // 50MB
		return nil, errors.New("file too large")
	}
	var data []byte
	data, err = ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}
	return &model.UploadFileReq{
		FileName: fileName,
		FileExt:  fileExt,
		FileSize: fileSize,
		FileData: data,
	}, err
}

func (d *UploadDecoder) UploadImage(_ context.Context, r *http.Request) (any, error) {
	f, fh, err := r.FormFile("file")
	if err != nil {
		return nil, err
	}
	defer f.Close()

	fileName, fileExt, fileSize, e := d.parseFileHeader(fh)
	if e != nil {
		return nil, e
	}
	if fileSize > 20*1024*1024 { // 20MB
		return nil, errors.New("image too large")
	}
	if matched, _ := regexp.MatchString("jpg|jpeg|png|gif|tif|tiff|bmp", fileExt); !matched {
		return nil, errors.New("unsupported image format")
	}
	var data []byte
	data, err = ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}
	compress := true
	if val := r.FormValue("compress"); val != "" {
		compress = gconv.Bool(val)
	}
	width := gconv.Int(r.FormValue("width"))
	height := gconv.Int(r.FormValue("height"))
	if width == 0 {
		width = 1000
	}
	if height == 0 {
		height = 1000
	}
	return &model.UploadImageReq{
		FileName: fileName,
		FileExt:  fileExt,
		FileSize: fileSize,
		FileData: data,
		Compress: compress,
		Width:    width,
		Height:   height,
	}, err
}

func (*UploadDecoder) parseFileHeader(fileHeader *multipart.FileHeader) (fileName, fileExt string, fileSize int64, err error) {
	if fileHeader == nil {
		err = errors.New("incorrect file")
		return
	}
	fileName = fileHeader.Filename
	fileExt = path.Ext(fileName)
	fileSize = fileHeader.Size
	return
}
