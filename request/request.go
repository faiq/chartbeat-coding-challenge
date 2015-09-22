package request

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
func MakeRequest(url string, updates chan<- Page) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	dec := json.NewDecoder(resp.Body)
	defer resp.Body.Close()
	//read open bracket
	_, err = dec.Token()
	if err != nil {
		return err
	}
	for dec.More() {
		var pageInfo Page
		err := dec.Decode(&pageInfo)
		if err != nil {
			return err
		}
		fmt.Printf("%v\n", pageInfo)
		updates <- pageInfo
	}
	// read closing bracket
	_, err = dec.Token()
	if err != nil {
		return err
	}
	return nil
}
