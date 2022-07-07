package endpoint

import (
	"elf-server/pkg/module/poststatistics/service"

	"github.com/mingslife/bone"
)

type PostStatisticsEndpoint struct {
	Service *service.PostStatisticsService `inject:""`
}

var _ bone.Endpoint = (*PostStatisticsEndpoint)(nil)
