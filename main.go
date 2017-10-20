package main

import (
	"fmt"
	"net/http"

	"google.golang.org/appengine"
)

func main() {
	http.HandleFunc("/", helloHandler)
	appengine.Main()
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Coolest 200 OK you've ever seen")
}
