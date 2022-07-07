package transport

import (
	"context"
	"errors"
	"net/http"
	"regexp"

	"elf-server/pkg/conf"
	"elf-server/pkg/module/log/model"
)

type LogDecoder struct{}

func (*LogDecoder) List(_ context.Context, r *http.Request) (any, error) {
	if !conf.GetConfig().Log {
		return nil, errors.New("log disabled")
	}

	month := r.URL.Query().Get("month")
	if ok, err := regexp.Match("^[0-9]{4}-[0-9]{2}$", []byte(month)); !ok || err != nil {
		return nil, errors.New("incorrect month")
	}

	return &model.ListReq{
		Month: month,
	}, nil
}

func (*LogDecoder) Get(_ context.Context, r *http.Request) (any, error) {
	if !conf.GetConfig().Log {
		return nil, errors.New("log disabled")
	}

	date := r.URL.Query().Get("date")
	if ok, err := regexp.Match("^[0-9]{4}-[0-9]{2}-[0-9]{2}$", []byte(date)); !ok || err != nil {
		return nil, errors.New("incorrect date")
	}

	return &model.GetReq{
		Date: date,
	}, nil
}
