package adapters

import "github.com/go-chi/chi/v5"

func InitAdapters(router *chi.Mux) {
	router.Route("/adapters", func(r chi.Router) {
		r.Route("/waha", func(r chi.Router) { initWaha(r) })
	})
}
