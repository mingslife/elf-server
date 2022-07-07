package setting

import (
	"github.com/mingslife/bone"

	"elf-server/pkg/module/setting/transport"
)

type Module struct {
	Http *transport.SettingHttp `inject:""`
}

func (*Module) Name() string {
	return "module.setting"
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
