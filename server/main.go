package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

type Cotacao struct {
	Code        string
	Codein      string
	Name        string
	High        string
	Low         string
	VarBid      string
	PctChange   string
	Bid         string
	Ask         string
	Timestamp   string
	Create_date string
}

type Dolar struct {
	Valor string
}

func main() {
	http.HandleFunc("/cotacao", handler)
	http.ListenAndServe(":8080", nil)
}

func handler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log.Println("\n\nRequest iniciada...")
	defer log.Println("Requst finalizada")
	select {
	case <-time.After(time.Second * 10):

		log.Println("Request processada com sucesso")

		dolar, err := buscarDadosApi()

		if err != nil {
			log.Println(err)
			panic(err)
		}
		//		err = salvar(dados)

		fmt.Println("teste", dolar.Valor)
		w.Write([]byte("Process ok"))
	case <-ctx.Done():
		log.Println("Request cancelada pele cliente")
	}

}

func buscarDadosApi() (*Dolar, error) {

	url := "https://economia.awesomeapi.com.br/json/last/USD-BRL"

	res, err := http.Get(url)
	if err != nil {
		log.Fatalf("Erro ao criar a requisição: %v", err)
		return nil, err
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		log.Fatalf("Erro ao fazer a requisição: %v", err)
		return nil, err
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatalf("Erro ao ler o corpo da resposta: %v", err)
		return nil, err
	}
	log.Println(body)
	var cotacao Cotacao
	err = json.Unmarshal(body, &cotacao)
	if err != nil {
		log.Fatalf("Erro ao fazer o parse do JSON: %v", err)
		return nil, err
	}

	log.Printf("teste 01: ", cotacao)

	return &Dolar{Valor: cotacao.Bid}, nil
}
