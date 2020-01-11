package router

import (
	"github.com/go-chi/chi"
	"gnodivad/sortify/src/api"
	"net/http"
)

func Init() http.Handler {
	r := chi.NewRouter()

	r.Get("/login", api.StartAuth)
	r.Get("/callback", api.CompleteAuth)

	return r
}
