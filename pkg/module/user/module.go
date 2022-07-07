package user

import (
	"github.com/mingslife/bone"

	"elf-server/pkg/module/user/transport"
)

type Module struct {
	Http *transport.UserHttp `inject:""`
}

func (*Module) Name() string {
	return "module.user"
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
