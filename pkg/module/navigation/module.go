package navigation

import (
	"github.com/mingslife/bone"

	"elf-server/pkg/module/navigation/transport"
)

type Module struct {
	Http *transport.NavigationHttp `inject:""`
}

func (*Module) Name() string {
	return "module.navigation"
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
