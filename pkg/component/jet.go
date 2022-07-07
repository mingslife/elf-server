package component

import (
	"context"
	"errors"
	"net/http"

	"github.com/CloudyKit/jet/v6"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/mingslife/bone"

	"elf-server/pkg/conf"
)

type PageData interface {
	TemplatePath() string
}

type Jet struct {
	jetSet *jet.Set
}

func (*Jet) Name() string {
	return "component.jet"
}

func (*Jet) Init() error {
	return nil
}

func (c *Jet) Register() error {
	cfg := conf.GetConfig()
	templateDir := "./templates"
	if cfg.Debug {
		c.jetSet = jet.NewSet(
			jet.NewOSFileSystemLoader(templateDir),
			jet.InDevelopmentMode(),
		)
	} else {
		c.jetSet = jet.NewSet(
			jet.NewOSFileSystemLoader(templateDir),
		)
	}
	return nil
}

func (*Jet) Unregister() error {
	return nil
}

func (c *Jet) RenderFunc(name string) httptransport.EncodeResponseFunc {
	t, err := c.jetSet.GetTemplate(name)
	if err != nil {
		panic(err)
	}
	return func(ctx context.Context, w http.ResponseWriter, rsp any) error {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		return t.Execute(w, nil, rsp)
	}
}

func (c *Jet) EncodeHTMLResponse(ctx context.Context, w http.ResponseWriter, rsp any) error {
	if v, ok := rsp.(PageData); ok {
		t, err := c.jetSet.GetTemplate(v.TemplatePath())
		if err != nil {
			return err
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		return t.Execute(w, nil, v)
	}
	return errors.New("incorrect response")
}

var _ bone.Component = (*Jet)(nil)
