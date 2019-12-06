package router

import (
	"github.com/go-chi/chi"
	"net/http"
)

func Init() http.Handler {
	r := chi.NewRouter()

	return r
}
