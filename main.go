package main

import (
	"log"
	"net/http"
	"task/handlers"
	"github.com/julienschmidt/httprouter"
)


func IndexHadler(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	w.WriteHeader(200)
	w.Write([]byte("Hello"))
}

func main() {

	router := httprouter.New()
	router.GET("/", IndexHadler)
	router.GET("/api/Create/", handlers.GenerateTokensHandler)
	router.GET("/api/Refresh/", handlers.RefreshTokensHandler)

	log.Fatal(http.ListenAndServe(":8080", router))
}
