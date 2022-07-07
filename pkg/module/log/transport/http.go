package transport

import (
	"context"
	"errors"
	"net/http"

	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/mingslife/bone"

	"elf-server/pkg/module/log/endpoint"
	"elf-server/pkg/module/log/model"
)

type LogHttp struct {
	Router   *bone.Router          `inject:"application.router"`
	Endpoint *endpoint.LogEndpoint `inject:""`
	Decoder  *LogDecoder           `inject:""`
}

func (t *LogHttp) Register() error {
	o, j := http.MethodOptions, httptransport.EncodeJSONResponse
	s := t.Router.PathPrefix("/api/v1/log").Subrouter()
	s.Methods(http.MethodGet, o).Path("/dates").Handler(httptransport.NewServer(t.Endpoint.List, t.Decoder.List, j))
	s.Methods(http.MethodGet, o).Path("/raw").Handler(httptransport.NewServer(t.Endpoint.Get, t.Decoder.Get, t.logResponse))
	return nil
}

func (t *LogHttp) logResponse(_ context.Context, w http.ResponseWriter, response any) error {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	if v, ok := response.(*model.GetRsp); ok {
		w.WriteHeader(http.StatusOK)
		_, err := w.Write(v.Data)
		return err
	}
	return errors.New("log not found")
}

var _ bone.Transport = (*LogHttp)(nil)
