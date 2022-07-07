package post

import (
	"github.com/mingslife/bone"

	"elf-server/pkg/module/post/transport"
)

type Module struct {
	Http *transport.PostHttp `inject:""`
}

func (*Module) Name() string {
	return "module.post"
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
