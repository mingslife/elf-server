package reader

import (
	"github.com/mingslife/bone"

	"elf-server/pkg/module/reader/transport"
)

type Module struct {
	Http *transport.ReaderHttp `inject:""`
}

func (*Module) Name() string {
	return "module.reader"
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
