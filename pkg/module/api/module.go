package api

import (
	"github.com/mingslife/bone"

	"elf-server/pkg/module/api/transport"
)

type Module struct {
	Http *transport.ApiHttp `inject:""`
}

func (*Module) Name() string {
	return "module.api"
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
