package transport

import (
	"net/http"

	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/mingslife/bone"

	"elf-server/pkg/module/upload/endpoint"
)

type UploadHttp struct {
	Router   *bone.Router             `inject:"application.router"`
	Endpoint *endpoint.UploadEndpoint `inject:""`
	Decoder  *UploadDecoder           `inject:""`
}

func (t *UploadHttp) Register() error {
	o, j := http.MethodOptions, httptransport.EncodeJSONResponse
	s := t.Router.PathPrefix("/api/v1/upload").Subrouter()
	s.Methods(http.MethodPost, o).Path("/file").Handler(httptransport.NewServer(t.Endpoint.UploadFile, t.Decoder.UploadFile, j))
	s.Methods(http.MethodPost, o).Path("/image").Handler(httptransport.NewServer(t.Endpoint.UploadImage, t.Decoder.UploadImage, j))
	return nil
}

var _ (bone.Transport) = (*UploadHttp)(nil)
