package main

import "net/http"

func main() {
	r, err := http.Post("http://letsboot-backend/schedule", "", nil)
	if err != nil {
		panic(err)
	}
	if r.StatusCode >= 400 || r.StatusCode < 200 {
		panic(r.Status)
	}
}
