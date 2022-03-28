package main

import (
	"golang-disributed-systems/internal/server"
	"log"
)

func main() {
	srv := server.NewHTTPServer(":8080")
	log.Printf("Server Running...")
	log.Fatal(srv.ListenAndServe())
}
