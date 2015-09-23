package main

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/codegangsta/negroni"
	"github.com/faiq/chartbeat-coding-challenge/request"
	"github.com/gorilla/mux"
)

const (
	pollInterval = 5 * time.Second                                                                   // how often to poll each URL
	baseUrl      = "http://api.chartbeat.com/live/toppages/?apikey=317a25eccba186e0f6b558f45214c0e7" //"Base Url" that we will make requests to chartbeat from
)

var state = make(map[string][]*PagePlus) // state will keep a mapping of hosts to corresponding pageplus structs
var mutex = &sync.Mutex{}                //keeps state in check between threads

//Page Plus is a struct that holds the same data along with a field that holds previous Visitors
type PagePlus struct {
	I            string `json:"i"`
	Path         string `json:"path"`
	Visitors     int    `json:"visitors"`
	PrevVisitors int
}

func MainHandler(w http.ResponseWriter, r *http.Request) {
	host := r.URL.Query().Get("host")
	updates := Poll(host)
	go func() {
		for {
			select {
			case page := <-updates:
				HandlePage(page, host)
			}
		}
	}()
}

// Poll will create a new channel for processing on a given host name every Interval seconds
func Poll(host string) chan request.Page {
	ticker := time.NewTicker(pollInterval)
	updates := make(chan request.Page) // a channel to pass along updates to
	hostString := baseUrl + "&host=" + host + "&limit=100"
	fmt.Printf(hostString)
	go func() {
		for {
			select {
			case <-ticker.C:
				err := request.MakeRequest(hostString, updates)
				if err != nil {
					fmt.Printf("%v", err)
					break
				}
			}
		}
	}()
	return updates
}

// handle page will take an incoming page struct and atomically save it to our global map
func HandlePage(page request.Page, host string) {
	mutex.Lock()
	if state[host] == nil {
		state[host] = append(state[host], &PagePlus{page.I, page.Path, page.Visitors, 0})
	}
	newPath := true // flag to determine whether or not to add a new member to the slice
	// loop over the saved paths in our state object
	for savedPage := range state[host] {
		if savedPage.Path == page.Path {
			newPath = false
			savedPage.PrevVisitors = savedPage.Visitors
			savedPage.Visitors = page.Visitors
		}
	}
	if newPath == true {
		state[host] = append(state[host], &PagePlus{page.I, page.Path, page.Visitors, 0})
	}
	mutex.Unlock()
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/", MainHandler)
	n := negroni.Classic()
	n.UseHandler(router)
	n.Run(":3000")
}
