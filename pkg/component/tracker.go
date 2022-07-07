package component

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	"github.com/mingslife/bone"
)

const (
	TraceIDKey = "X-Trace-Id"
)

type Tracker struct {
	Router *bone.Router `inject:"application.router"`
}

func (*Tracker) Name() string {
	return "component.tracker"
}

func (*Tracker) Init() error {
	return nil
}

func (c *Tracker) Register() error {
	c.Router.Use(c.Middleware)
	return nil
}

func (*Tracker) Unregister() error {
	return nil
}

func (*Tracker) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		traceID := uuid.New().String()
		ctx := context.WithValue(r.Context(), TraceIDKey, traceID)
		*r = *r.Clone(ctx)
		r.Header.Add(TraceIDKey, traceID)
		w.Header().Set(TraceIDKey, traceID)

		next.ServeHTTP(w, r)
	})
}

var _ bone.Component = (*Tracker)(nil)
