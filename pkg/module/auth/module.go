package auth

import (
	"github.com/mingslife/bone"

	"elf-server/pkg/module/auth/transport"
)

type Module struct {
	Http *transport.AuthHttp `inject:""`
}

func (*Module) Name() string {
	return "module.auth"
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
