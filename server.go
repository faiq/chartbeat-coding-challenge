package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/codegangsta/negroni"
	"github.com/faiq/chartbeat-coding-challenge/request"
	"github.com/garyburd/redigo/redis"
	"github.com/gorilla/mux"
)

const (
	pollInterval = 5 * time.Second                                                                   // how often to poll each URL
	baseUrl      = "http://api.chartbeat.com/live/toppages/?apikey=317a25eccba186e0f6b558f45214c0e7" //"Base Url" that we will make requests to chartbeat from
)

var pool = newPool()

func MainHandler(w http.ResponseWriter, r *http.Request) {
	host := r.URL.Query().Get("host")
	updates := Poll(host)
	go func() {
		for {
			select {
			case page, chanClosed := <-updates:
				if !(chanClosed) {
					fmt.Printf("Channel Closed")
				} else {
					fmt.Printf("this is page %v", page)
					HandlePage(page)
				}
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

// handle page will take an incoming page struct and do some processing on it
func HandlePage(page request.Page) {
	c := pool.Get()
	defer c.Close()
	str := page.Path
	prevNum, err := redis.Int(c.Do("GET", str))
	fmt.Printf("this is prevNum %d \n", prevNum)
	if err != nil && err != redis.ErrNil {
		fmt.Printf("Redis is throwing an error getting this key %s and this is the err %v", str, err)
	}
	if prevNum < page.Visitors {
		diff := page.Visitors - prevNum
		_, err := c.Do("SET", str, diff)
		if err != nil {
			fmt.Printf("Error writing to Redis %v \n", err)
		}
	}
}

func newPool() *redis.Pool {
	return &redis.Pool{
		MaxIdle:   80,
		MaxActive: 12000, // max number of connections
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", ":6379")
			if err != nil {
				panic(err.Error())
			}
			return c, err
		},
	}

}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/", MainHandler)
	n := negroni.Classic()
	n.UseHandler(router)
	n.Run(":3000")
}
