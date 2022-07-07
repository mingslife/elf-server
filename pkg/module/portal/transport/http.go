package transport

import (
	"net/http"

	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/mingslife/bone"

	"elf-server/pkg/component"
	"elf-server/pkg/module/portal/endpoint"
)

type PortalHttp struct {
	Router   *bone.Router             `inject:"application.router"`
	Jet      *component.Jet           `inject:"component.jet"`
	Endpoint *endpoint.PortalEndpoint `inject:""`
	Decoder  *PortalDecoder           `inject:""`
}

func (t *PortalHttp) Register() error {
	h := t.Jet.EncodeHTMLResponse
	s := t.Router.NewRoute().Subrouter()
	s.Methods(http.MethodGet).Path("/").Handler(httptransport.NewServer(t.Endpoint.Index, t.Decoder.Index, h))
	s.Methods(http.MethodGet).Path("/post/{route}").Handler(httptransport.NewServer(t.Endpoint.Post, t.Decoder.Post, h))
	s.Methods(http.MethodGet).Path("/user/{username}").Handler(httptransport.NewServer(t.Endpoint.User, t.Decoder.User, h))
	s.Methods(http.MethodGet).Path("/user/{username}/{page}").Handler(httptransport.NewServer(t.Endpoint.User, t.Decoder.User, h))
	s.Methods(http.MethodGet).Path("/category/{route}").Handler(httptransport.NewServer(t.Endpoint.Category, t.Decoder.Category, h))
	s.Methods(http.MethodGet).Path("/category/{route}/{page}").Handler(httptransport.NewServer(t.Endpoint.Category, t.Decoder.Category, h))
	s.Methods(http.MethodGet).Path("/posts").Handler(httptransport.NewServer(t.Endpoint.Posts, t.Decoder.Posts, h))
	s.Methods(http.MethodGet).Path("/posts/{page}").Handler(httptransport.NewServer(t.Endpoint.Posts, t.Decoder.Posts, h))
	s.Methods(http.MethodGet).Path("/reader").Handler(httptransport.NewServer(t.Endpoint.Reader, t.Decoder.Reader, h))
	s.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", http.FileServer(http.Dir("./assets"))))
	s.PathPrefix("/upload/").Handler(http.StripPrefix("/upload/", http.FileServer(http.Dir("./upload"))))
	return nil
}

var _ bone.Transport = (*PortalHttp)(nil)
