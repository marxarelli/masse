package main

import (
	"fmt"
	"io"
	"net/http"

	"nmyk.io/cowsay"
)

func handler(w http.ResponseWriter, r *http.Request) {
	cow := cowsay.Cow{}

	moo, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "the cow doth protest", http.StatusBadRequest)
		return
	}

	cow.Write(w, moo, false)
}

func main() {
	http.HandleFunc("/", handler)
	fmt.Println("cowsayd is running on port 3000")
	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}
