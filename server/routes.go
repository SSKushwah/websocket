package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

// SetupRoutes provides all the routes that can be used
func (srv *Server) InjectRoutes() *chi.Mux {
	router := chi.NewRouter()
	// router.Route("/v1", func(v1 chi.Router) {

	router.Route("/user", func(user chi.Router) {
		// user.Use(middlewares.AuthMiddleware)

		user.Post("/check", func(w http.ResponseWriter, r *http.Request) {

			// fmt.Fprintf(w, "Hello, you've requested: %s\n", r.URL.Path)

			str := "test is checked"

			w.Write([]byte(str))
		})

		user.HandleFunc("/connect", srv.Connect)

	})
	// })
	return router
}
