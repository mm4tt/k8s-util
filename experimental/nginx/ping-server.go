package main

import (
	"encoding/json"
	"fmt"
	"html"
	"io/ioutil"
	"log"
	"net/http"
)

func main() {
	fmt.Println("I'm ping server.")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		body, _ := ioutil.ReadAll(r.Body)
		headers, _ := json.MarshalIndent(r.Header, "", "  ")

		fmt.Fprintf(w, "URL: %s\n", html.EscapeString(r.URL.String()))
		fmt.Fprintf(w, "Body: %q\n", string(body))
		fmt.Fprintf(w, "Headers: %s", string(headers))
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
