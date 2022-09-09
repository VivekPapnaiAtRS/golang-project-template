package server

import (
	"github.com/go-chi/chi/v5"
)

func (srv *Server) InjectRoutes() *chi.Mux {
	router := chi.NewRouter()
	router.Route("/api", func(api chi.Router) {
		api.Get("/welcome", srv.greet)
		api.Get("/ws", srv.realTimeWS)
	})
	return router
}
