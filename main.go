package main

import (
	"fmt"
	"net/http"

	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
)

func main() {
	http.HandleFunc("/", helloHandler)
	http.HandleFunc("/sigfox/callback", handleCallback)
	appengine.Main()
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Coolest 200 OK you've ever seen")
}

func interalServerError(w http.ResponseWriter, r *http.Request, err error) {
	ctx := appengine.NewContext(r)
	log.Debugf(ctx, "Internal server error: %v", err)
	http.Error(w, "Internal server error.", http.StatusInternalServerError)
}
