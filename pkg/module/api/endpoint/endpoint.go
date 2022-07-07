package endpoint

import (
	"bytes"
	"context"
	"elf-server/pkg/module/api/model"

	"github.com/dchest/captcha"
	"github.com/mingslife/bone"
)

type ApiEndpoint struct{}

func (*ApiEndpoint) Status(ctx context.Context, req any) (any, error) {
	return nil, nil
}

func (*ApiEndpoint) Captcha(ctx context.Context, req any) (any, error) {
	var buf bytes.Buffer

	captchaID := captcha.New()
	captcha.WriteImage(&buf, captchaID, captcha.StdWidth, captcha.StdHeight)

	data := buf.Bytes()

	return &model.CaptchaRsp{
		CaptchaID: captchaID,
		Data:      data,
	}, nil
}

var _ bone.Endpoint = (*ApiEndpoint)(nil)
