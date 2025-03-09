package main

import (
	"net/http"
)

func main() {
	http.Handle("/gotoproto/", http.StripPrefix("/gotoproto/", http.FileServer(http.Dir("./docs"))))
	http.ListenAndServe(":3000", nil)
}
