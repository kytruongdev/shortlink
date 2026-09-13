package handler

import (
	"github.com/go-chi/chi/v5"

	"github.com/kytruongdev/shortlink/internal/handler/rest"
	"github.com/kytruongdev/shortlink/internal/infra/httpserver"
)

// Router holds business route dependencies.
type Router struct {
	restHandler *rest.Handler
	authHandler *rest.AuthHandler
	jwtSecret   []byte
}

// New returns a Router wiring the link and auth REST handlers.
func New(restHandler *rest.Handler, authHandler *rest.AuthHandler, jwtSecret []byte) Router {
	return Router{restHandler: restHandler, authHandler: authHandler, jwtSecret: jwtSecret}
}

// Routes registers the service routes on r. optionalAuth attaches the user (if any)
// to every /api/v1 request; the authenticated group additionally enforces login.
func (rtr Router) Routes(r chi.Router) {
	r.Route("/api/v1", func(r chi.Router) {
		r.Use(optionalAuth(rtr.jwtSecret))
		r.Group(rtr.public)
		r.Group(rtr.authenticated)
	})
}

// public registers routes that work with or without a logged-in user.
func (rtr Router) public(r chi.Router) {
	r.Post("/encode", httpserver.HandlerErr(rtr.restHandler.Encode))
	r.Post("/decode", httpserver.HandlerErr(rtr.restHandler.Decode))

	// Frontend resolves a short code here, then performs the client-side redirect.
	r.Get("/{code}", httpserver.HandlerErr(rtr.restHandler.Resolve))

	r.Route("/auth", func(r chi.Router) {
		r.Post("/register", httpserver.HandlerErr(rtr.authHandler.Register))
		r.Post("/login", httpserver.HandlerErr(rtr.authHandler.Login))
		r.Post("/refresh", httpserver.HandlerErr(rtr.authHandler.Refresh))
		r.Post("/logout", httpserver.HandlerErr(rtr.authHandler.Logout))
	})
}

// authenticated registers routes that require a logged-in user.
func (rtr Router) authenticated(r chi.Router) {
	r.Use(requireAuth)
	r.Get("/links", httpserver.HandlerErr(rtr.restHandler.ListLinks))
	r.Patch("/links/{code}", httpserver.HandlerErr(rtr.restHandler.UpdateLink))
	r.Delete("/links/{code}", httpserver.HandlerErr(rtr.restHandler.DeleteLink))
}
