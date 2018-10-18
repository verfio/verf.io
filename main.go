package main

import (
	"log"
	"net/http"
)

func main() {

	//fmt.Println(t)
	//r := mux.NewRouter()
	srv := http.NewServeMux()
	srv.Handle("/", http.FileServer(http.Dir(".")))

	log.Printf("server started")

	log.Fatal(http.ListenAndServe(":3333", srv))

}
