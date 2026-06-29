package handler

import (
	"github.com/go-chi/chi/v5"

	"github.com/kytruongdev/shortlink/internal/handler/rest"
	"github.com/kytruongdev/shortlink/internal/infra/httpserver"
)

// Router holds business route dependencies.
type Router struct {
	restHandler *rest.Handler
}

// New returns a Router wiring the given REST handler.
func New(restHandler *rest.Handler) Router {
	return Router{restHandler: restHandler}
}

// Routes registers the service routes on r.
func (rtr Router) Routes(r chi.Router) {
	r.Group(rtr.public)
}

// public registers routes that require no authentication.
func (rtr Router) public(r chi.Router) {
	r.Route("/api/v1", func(r chi.Router) {
		r.Post("/encode", httpserver.HandlerErr(rtr.restHandler.Encode))
		r.Post("/decode", httpserver.HandlerErr(rtr.restHandler.Decode))
	})
}
