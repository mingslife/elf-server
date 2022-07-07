package log

import (
	"github.com/mingslife/bone"

	"elf-server/pkg/module/log/transport"
)

type Module struct {
	Http *transport.LogHttp `inject:""`
}

func (*Module) Name() string {
	return "module.log"
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
