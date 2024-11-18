package main

import (
	"net/http"
)

func main() {
	http.HandleFunc("/cotacao", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Fazer chamada de api\n"))
	})
	http.ListenAndServe(":8080", nil)
}
