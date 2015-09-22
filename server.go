package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/faiq/chartbeat-coding-challenge/request"

	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
)

const (
	pollInterval = 5 * time.Second // how often to poll each URL
)

func MainHandler(w http.ResponseWriter, r *http.Request) {
	u := r.URL.Query()
	fmt.Printf("%v", u)
	w.Write([]byte("Gorilla!\n"))
}

func Poll(updates chan<- request.Page) {
	ticker := time.NewTicker(updateInterval)
	go func() {
		for {
			select {
			case <-ticker.C:
				request.MakeRequest()
			}
		}
	}()
}

func main() {
	router := mux.NewRouter()
	updates := make(chan request.Page) // a channel to pass along updates to
	router.HandleFunc("/", MainHandler)
	n := negroni.Classic()
	n.UseHandler(router)
	n.Run(":3000")
}
