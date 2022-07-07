package transport

import (
	"net/http"

	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/mingslife/bone"

	"elf-server/pkg/module/auth/endpoint"
)

type AuthHttp struct {
	Router   *bone.Router           `inject:"application.router"`
	Endpoint *endpoint.AuthEndpoint `inject:""`
	Decoder  *AuthDecoder           `inject:""`
}

func (t *AuthHttp) Register() error {
	o, j := http.MethodOptions, httptransport.EncodeJSONResponse
	s := t.Router.PathPrefix("/api/v1/auth").Subrouter()
	s.Methods(http.MethodPost, o).Path("/login").Handler(httptransport.NewServer(t.Endpoint.Login, t.Decoder.Login, j))
	s.Methods(http.MethodPost, o).Path("/logout").Handler(httptransport.NewServer(t.Endpoint.Logout, t.Decoder.Logout, j))
	s.Methods(http.MethodPost, o).Path("/refresh").Handler(httptransport.NewServer(t.Endpoint.Refresh, t.Decoder.Refresh, j))
	s.Methods(http.MethodGet, o).Path("/settings").Handler(httptransport.NewServer(t.Endpoint.GetSettings, t.Decoder.GetSettings, j))
	s.Methods(http.MethodGet, o).Path("/info").Handler(httptransport.NewServer(t.Endpoint.GetInfo, t.Decoder.GetInfo, j))
	return nil
}

var _ bone.Transport = (*AuthHttp)(nil)
