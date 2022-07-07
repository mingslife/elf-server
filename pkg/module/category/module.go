package category

import (
	"github.com/mingslife/bone"

	"elf-server/pkg/module/category/transport"
)

type Module struct {
	Http *transport.CategoryHttp `inject:""`
}

func (*Module) Name() string {
	return "module.category"
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
