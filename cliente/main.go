package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, "GET", "http://localhost:8080/cotacao", nil)
	if err != nil {
		panic(err)
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		fmt.Println("Erro no resultaldo")
		return
	}

	file, err := os.Create("arquivo.txt")
	if err != nil {
		fmt.Println("Erro ao criar arquivo")
		return
	}

	defer file.Close()

	_, err = io.Copy(file, res.Body)
	if err != nil {
		fmt.Println("Erro ao escrever arquivo")
		return
	}
	log.Print("Processo finalizado")
}
