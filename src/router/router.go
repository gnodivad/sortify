package router

import (
	"gnodivad/sortify/src/api"
	"net/http"

	"github.com/go-chi/chi"
)

func Init() http.Handler {
	r := chi.NewRouter()

	r.Get("/login", api.StartAuth)
	r.Get("/callback", api.CompleteAuth)
	r.Get("/welcome", api.Welcome)

	return r
}
