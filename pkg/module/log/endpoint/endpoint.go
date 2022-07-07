package endpoint

import (
	"context"
	"elf-server/pkg/component"
	"elf-server/pkg/module/log/model"
	"elf-server/pkg/utils"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
)

type LogEndpoint struct{}

func (*LogEndpoint) List(ctx context.Context, req any) (any, error) {
	r := req.(*model.ListReq)
	var dates []string
	path := component.LogPath()
	for day := 1; day <= 31; day++ {
		dayStr := strconv.Itoa(day)
		if day < 10 {
			dayStr = "0" + dayStr
		}

		ymd := r.Month + "-" + dayStr
		fileName := component.LogFileName(ymd)
		filePath := filepath.Join(path, fileName)
		if utils.IsFileExists(filePath) {
			dates = append(dates, ymd)
		}
	}
	return &model.ListRsp{
		Dates: dates,
	}, nil
}

func (*LogEndpoint) Get(ctx context.Context, req any) (any, error) {
	r := req.(*model.GetReq)
	path := component.LogPath()
	ymd := r.Date
	fileName := component.LogFileName(ymd)
	filePath := filepath.Join(path, fileName)
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var data []byte
	data, err = ioutil.ReadAll(f)
	return &model.GetRsp{
		Data: data,
	}, err
}
