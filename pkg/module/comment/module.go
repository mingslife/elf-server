package comment

import (
	"github.com/mingslife/bone"

	"elf-server/pkg/module/comment/transport"
)

type Module struct {
	Http *transport.CommentHttp `inject:""`
}

func (*Module) Name() string {
	return "module.comment"
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
