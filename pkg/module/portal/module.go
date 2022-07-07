package portal

import (
	"github.com/mingslife/bone"

	"elf-server/pkg/module/portal/transport"
)

type Module struct {
	Http *transport.PortalHttp `inject:""`
}

func (*Module) Name() string {
	return "module.portal"
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
