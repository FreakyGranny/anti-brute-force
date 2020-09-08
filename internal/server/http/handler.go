package internalhttp

import (
	"io"
	"net/http"

	"github.com/FreakyGranny/anti-brute-force/internal/app"
)

// HealthcheckHandler ...
type HealthcheckHandler struct {
	App app.Application
}

// NewHealthcheckHandler ...
func NewHealthcheckHandler(a app.Application) *HealthcheckHandler {
	return &HealthcheckHandler{App: a}
}

func (h *HealthcheckHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/healthcheck":
		h.HealthCheck(w, r)
	default:
		http.NotFound(w, r)
	}
}

// HealthCheck simple route.
func (h *HealthcheckHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "OK") //nolint
}
