package transport

import (
	"github.com/mingslife/bone"

	"elf-server/pkg/module/poststatistics/endpoint"
)

type PostStatisticsHttp struct {
	Router   *bone.Router                     `inject:"application.router"`
	Endpoint *endpoint.PostStatisticsEndpoint `inject:""`
}

func (t *PostStatisticsHttp) Register() error {
	return nil
}

var _ bone.Transport = (*PostStatisticsHttp)(nil)
