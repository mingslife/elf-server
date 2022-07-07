package upload

import (
	"github.com/mingslife/bone"

	"elf-server/pkg/module/upload/transport"
)

type Module struct {
	Http *transport.UploadHttp `inject:""`
}

func (*Module) Name() string {
	return "module.upload"
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
