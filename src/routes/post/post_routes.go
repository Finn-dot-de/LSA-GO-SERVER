package post

import (
	"github.com/Finn-dot-de/LernStoffAnwendung/src/handler"
	"github.com/go-chi/chi"
)

// DefinePostRoutes definiert alle POST-Routen der Anwendung
func DefinePostRoutes(r *chi.Mux) {
	// Login-Handler
	r.Post("/api/login", handler.LoginHandler)
}
