package transport

import (
	"context"
	"net/http"
)

type ApiDecoder struct{}

func (*ApiDecoder) Status(_ context.Context, r *http.Request) (any, error) {
	return nil, nil
}

func (*ApiDecoder) Captcha(_ context.Context, r *http.Request) (any, error) {
	return nil, nil
}
