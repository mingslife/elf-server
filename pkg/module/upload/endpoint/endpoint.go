package endpoint

import (
	"context"
	"path/filepath"
	"strings"

	"github.com/mingslife/bone"

	"elf-server/pkg/module/upload/model"
	"elf-server/pkg/module/upload/service"
)

type UploadEndpoint struct {
	Service *service.UploadService `inject:""`
}

func (e *UploadEndpoint) UploadFile(ctx context.Context, req any) (any, error) {
	r := req.(*model.UploadFileReq)
	filePath, err := e.Service.SaveFile(r.FileName, r.FileData)
	if err != nil {
		return nil, err
	}
	return &model.UploadFileRsp{
		Path: e.getPath(filePath),
	}, nil
}

func (e *UploadEndpoint) UploadImage(ctx context.Context, req any) (any, error) {
	r := req.(*model.UploadImageReq)
	filePath, err := e.Service.SaveImage(r.FileName, r.FileData, r.Compress, r.Width, r.Height)
	if err != nil {
		return nil, err
	}
	return &model.UploadImageRsp{
		Path: e.getPath(filePath),
	}, nil
}

func (*UploadEndpoint) getPath(filePath string) string {
	path := strings.ReplaceAll(filepath.Clean("/"+filePath), `\`, `/`)
	return path
}

var _ bone.Endpoint = (*UploadEndpoint)(nil)
