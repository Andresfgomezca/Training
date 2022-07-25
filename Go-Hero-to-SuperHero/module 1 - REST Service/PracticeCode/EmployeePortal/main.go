package main

import (
	"log"
	"net/http"
	"time"
	"employeeportal/controller"
)

func main() {
	r := controller.Handlers()
	srv := http.Server{
		Handler:      r,
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 10,
		Addr:         "127.0.0.1:8001",
	}

	log.Println("Spinning up a new Mux Server...")

	log.Fatal(srv.ListenAndServe())
	//connection rejected to 8001
}
