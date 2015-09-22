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
	ticker := time.NewTicker(pollInterval)
	go func() {
		for {
			select {
			case <-ticker.C:
				fmt.Printf("tick tock")
				err := request.MakeRequest("http://api.chartbeat.com/live/toppages/?apikey=317a25eccba186e0f6b558f45214c0e7&host=gizmodo.com&limit=100", updates)
				if err != nil {
					fmt.Printf("Uh")
					break
				}
			}
		}
	}()
}

func main() {
	router := mux.NewRouter()
	updates := make(chan request.Page) // a channel to pass along updates to
	Poll(updates)
	go func() {
		for {
			select {
			case page := <-updates:
				fmt.Printf("ayy %v", page)
			}
		}
	}()
	router.HandleFunc("/", MainHandler)
	n := negroni.Classic()
	n.UseHandler(router)
	n.Run(":3000")
}
