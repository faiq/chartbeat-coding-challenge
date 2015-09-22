package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Page struct {
	I        string `json:"i"`
	Path     string `json:"path"`
	Visitors int    `json:"visitors"`
}

// Make request to the given url
func makeRequest(url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	dec := json.NewDecoder(resp.Body)
	//read open bracket
	_, err = dec.Token()
	if err != nil {
		return err
	}
	var pageInfo Page
	for dec.More() {
		err := dec.Decode(&pageInfo)
		if err != nil {
			return err
		}
		fmt.Printf("%v\n", pageInfo)
	}
	// read closing bracket
	_, err = dec.Token()
	if err != nil {
		return err
	}
	return nil
}

func main() {
	err := makeRequest("http://api.chartbeat.com/live/toppages/?apikey=317a25eccba186e0f6b558f45214c0e7&host=gizmodo.com&limit=100")
	if err != nil {
		fmt.Printf("JUMPMANNNN SHIT BROKE")
	}
}
