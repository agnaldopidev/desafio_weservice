package main

import (
	"context"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {

	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*300)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, "GET", "http://localhost:8080/cotacao", nil)
	if err != nil {
		log.Println(err.Error())
		return
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println(err.Error())
		return
	}
	defer res.Body.Close()

	if res.Body == nil {
		panic("Erro na respota")
	}

	if res.StatusCode != http.StatusOK {
		log.Println("Erro no resultaldo")
		return
	}

	file, err := os.Create("cotacao.txt")
	if err != nil {
		log.Println(err.Error())
	}

	defer file.Close()
	_, err = io.Copy(file, res.Body)
	if err != nil {
		log.Println(err.Error())
	}

	log.Print("Processo finalizado")
}
