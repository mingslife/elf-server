package component

import (
	"fmt"
	"net/http"

	"github.com/mingslife/bone"
)

type Error struct {
	Router *bone.Router `inject:"application.router"`
}

func (*Error) Name() string {
	return "component.error"
}

func (*Error) Init() error {
	return nil
}

func (c *Error) Register() error {
	c.Router.Use(c.Middleware)
	return nil
}

func (*Error) Unregister() error {
	return nil
}

func (*Error) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			err := recover()

			fmt.Println("error")
			w.Write([]byte("error"))

			switch err.(type) {
			case error:
			default:
			}
		}()

		next.ServeHTTP(w, r)
	})
}

var _ bone.Component = (*Error)(nil)
