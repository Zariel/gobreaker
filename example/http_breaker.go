package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/sony/gobreaker"
)

var cb *gobreaker.CircuitBreaker[[]byte]

func init() {
	var st gobreaker.Settings
	st.Name = "HTTP GET"
	st.ReadyToTrip = func(counts gobreaker.Counts) bool {
		failureRatio := float64(counts.TotalFailures) / float64(counts.Requests)
		return counts.Requests >= 3 && failureRatio >= 0.6
	}

	cb = gobreaker.NewCircuitBreaker[[]byte](st)
}

// Get wraps http.Get in CircuitBreaker.
func Get(url string) ([]byte, error) {
	return cb.Execute(func() ([]byte, error) {
		resp, err := http.Get(url)
		if err != nil {
			return nil, err
		}

		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		return body, nil
	})
}

func main() {
	body, err := Get("http://www.google.com/robots.txt")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(body))
}
