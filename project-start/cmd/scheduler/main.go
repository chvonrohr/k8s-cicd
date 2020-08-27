package main

import (
	"net/http"
	"os"
)

func main() {
	// default scheduler endpoint
	url := "http://letsboot-backend/schedule"

	// allow overriding of endpoint via argument
	if len(os.Args) > 1 {
		url = os.Args[1]
	}

	// run a post against the endpoint
	r, err := http.Post(url, "", nil)
	// if anything goes wrong, exit with a failed status => kubernetes will recognise the invalid exit code
	if err != nil {
		panic(err)
	}
	if r.StatusCode >= 400 || r.StatusCode < 200 {
		panic(r.Status)
	}
}
