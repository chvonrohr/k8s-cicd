package main

import (
	"net/http"
	"os"
)

func main() {
	url := "http://letsboot-backend/schedule"

	if len(os.Args) > 1 {
		url = os.Args[1]
	}

	r, err := http.Post(url, "", nil)
	if err != nil {
		panic(err)
	}
	if r.StatusCode >= 400 || r.StatusCode < 200 {
		panic(r.Status)
	}
}
