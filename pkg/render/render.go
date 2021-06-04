package render

import (
	"net/http"

	"github.com/CloudyKit/jet/v6"
	"github.com/gin-gonic/gin/render"
)

// see: https://github.com/amoniacou/ginjet

// RenderOptions is options's struct for JetRender
type RenderOptions struct {
	TemplateDir     string
	ContentType     string
	DevelopmentMode bool
}

// DefaultOptions is default options for JetRender
func DefaultOptions() *RenderOptions {
	return &RenderOptions{
		TemplateDir:     "./templates",
		ContentType:     "text/html; charset=utf-8",
		DevelopmentMode: true,
	}
}

// JetRender is a custom Gin template renderer using Jet
type JetRender struct {
	Options  *RenderOptions
	Template *jet.Template
	Data     interface{}
}

// New creates a new JetRender instance with custom Options.
func New(options *RenderOptions) *JetRender {
	return &JetRender{
		Options: options,
	}
}

// Default creates a JetRender instance with default options.
func Default() *JetRender {
	return New(DefaultOptions())
}

func (r JetRender) Instance(name string, data interface{}) render.Render {
	var set *jet.Set
	if r.Options.DevelopmentMode {
		set = jet.NewSet(
			jet.NewOSFileSystemLoader(r.Options.TemplateDir),
			jet.InDevelopmentMode(),
		)
	} else {
		set = jet.NewSet(
			jet.NewOSFileSystemLoader(r.Options.TemplateDir),
		)
	}
	t, err := set.GetTemplate(name)

	if err != nil {
		panic(err)
	}

	return JetRender{
		Data:     data,
		Options:  r.Options,
		Template: t,
	}
}

func (r JetRender) Render(w http.ResponseWriter) error {
	// Unless already set, write the Content-Type header.
	r.WriteContentType(w)
	if err := r.Template.Execute(w, nil, r.Data); err != nil {
		return err
	}
	return nil
}

func (r JetRender) WriteContentType(w http.ResponseWriter) {
	header := w.Header()
	if val := header["Content-Type"]; len(val) == 0 {
		header["Content-Type"] = []string{r.Options.ContentType}
	}
}
