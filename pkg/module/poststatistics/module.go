package poststatistics

import (
	"github.com/mingslife/bone"

	"elf-server/pkg/module/poststatistics/transport"
)

type Module struct {
	Http *transport.PostStatisticsHttp `inject:""`
}

func (*Module) Name() string {
	return "module.poststatistics"
}

func (*Module) Init() error {
	return nil
}

func (m *Module) Register() error {
	return m.Http.Register()
}

func (*Module) Unregister() error {
	return nil
}

var _ bone.Module = (*Module)(nil)
