package component

import (
	"net/http"

	"github.com/mingslife/bone"
)

type Cors struct {
	Router *bone.Router `inject:"application.router"`
}

func (*Cors) Name() string {
	return "component.cors"
}

func (*Cors) Init() error {
	return nil
}

func (c *Cors) Register() error {
	c.Router.Use(c.Middleware)
	return nil
}

func (*Cors) Unregister() error {
	return nil
}

func (*Cors) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "*")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		if r.Method == http.MethodOptions {
			return
		}

		next.ServeHTTP(w, r)
	})
}

var _ bone.Component = (*Cors)(nil)
