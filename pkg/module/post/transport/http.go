package transport

import (
	"net/http"

	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/mingslife/bone"

	"elf-server/pkg/module/post/endpoint"
)

type PostHttp struct {
	Router   *bone.Router           `inject:"application.router"`
	Endpoint *endpoint.PostEndpoint `inject:""`
	Decoder  *PostDecoder           `inject:""`
}

func (t *PostHttp) Register() error {
	o, j := http.MethodOptions, httptransport.EncodeJSONResponse
	s := t.Router.PathPrefix("/api/v1/posts").Subrouter()
	s.Methods(http.MethodPut, o).Path("/{id}/content").Handler(httptransport.NewServer(t.Endpoint.UpdateContent, t.Decoder.UpdateContent, j))
	s.Methods(http.MethodGet, o).Path("/{id}/content").Handler(httptransport.NewServer(t.Endpoint.GetContent, t.Decoder.GetContent, j))
	s.Methods(http.MethodGet, o).Path("").Handler(httptransport.NewServer(t.Endpoint.List, t.Decoder.List, j))
	s.Methods(http.MethodGet, o).Path("/{id}").Handler(httptransport.NewServer(t.Endpoint.Get, t.Decoder.Get, j))
	s.Methods(http.MethodPost, o).Path("").Handler(httptransport.NewServer(t.Endpoint.Create, t.Decoder.Create, j))
	s.Methods(http.MethodPut, o).Path("/{id}").Handler(httptransport.NewServer(t.Endpoint.Update, t.Decoder.Update, j))
	s.Methods(http.MethodDelete, o).Path("/{id}").Handler(httptransport.NewServer(t.Endpoint.Delete, t.Decoder.Delete, j))
	return nil
}

var _ bone.Transport = (*PostHttp)(nil)
