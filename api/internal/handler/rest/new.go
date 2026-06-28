package rest

import (
	controllerlink "github.com/kytruongdev/shortlink/internal/controller/link"
)

// Handler serves the v1 REST endpoints.
type Handler struct {
	ctrl    controllerlink.Controller
	baseURL string
}

// New returns a REST Handler.
func New(ctrl controllerlink.Controller, baseURL string) *Handler {
	return &Handler{ctrl: ctrl, baseURL: baseURL}
}
