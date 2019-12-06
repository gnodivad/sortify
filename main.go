package main

import (
	"gnodivad/sortify/src/router"
	"net/http"
)

func main() {
	r := router.Init()
	http.ListenAndServe(":8080", r)
}
