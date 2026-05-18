package main

import (
	"fmt"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
              fmt.Fprintln(w, "mini-ai-compute-platform is running")
	})

	fmt.Println("server started at :8080")

	if err := http.ListenAndServe(":8080", mux); err != nil {
		panic(err)
	}
}
