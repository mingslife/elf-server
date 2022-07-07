package transport

import (
	"context"
	"errors"
	"net/http"
	"net/url"

	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/mingslife/bone"

	"elf-server/pkg/module/api/endpoint"
	"elf-server/pkg/module/api/model"
)

type ApiHttp struct {
	Router   *bone.Router          `inject:"application.router"`
	Endpoint *endpoint.ApiEndpoint `inject:""`
	Decoder  *ApiDecoder           `inject:""`
}

func (t *ApiHttp) Register() error {
	o, j := http.MethodOptions, httptransport.EncodeJSONResponse
	s := t.Router.NewRoute().Subrouter()
	s.Methods(http.MethodGet, o).Path("/content/{uniqueId}").Handler(httptransport.NewServer(t.Endpoint.Status, t.Decoder.Status, j))
	s.Methods(http.MethodGet, o).Path("/comment/{uniqueId}").Handler(httptransport.NewServer(t.Endpoint.Status, t.Decoder.Status, j))
	s.Methods(http.MethodPost, o).Path("/comment/{uniqueId}").Handler(httptransport.NewServer(t.Endpoint.Status, t.Decoder.Status, j))
	s.Methods(http.MethodGet, o).Path("/captcha").Handler(httptransport.NewServer(t.Endpoint.Captcha, t.Decoder.Captcha, t.captchaResponse))
	s.Methods(http.MethodPost, o).Path("/reader/login")
	s.Methods(http.MethodGet, o).Path("/reader/info").Handler(httptransport.NewServer(t.Endpoint.Status, t.Decoder.Status, j))
	return nil
}

func (t *ApiHttp) captchaResponse(_ context.Context, w http.ResponseWriter, response any) error {
	w.Header().Set("Content-Type", "image/png")
	if v, ok := response.(*model.CaptchaRsp); ok {
		http.SetCookie(w, &http.Cookie{
			Name:     "ELF_CAPTCHA_ID",
			Value:    url.QueryEscape(v.CaptchaID),
			MaxAge:   600,
			Path:     "/",
			Domain:   "",
			SameSite: http.SameSiteDefaultMode,
			Secure:   false,
			HttpOnly: true,
		})
		w.WriteHeader(http.StatusOK)
		_, err := w.Write(v.Data)
		return err
	}
	return errors.New("captcha not found")
}

var _ bone.Transport = (*ApiHttp)(nil)
