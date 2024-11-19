package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type Dados struct {
	USDBRL Cotacao `json:"USDBRL"`
}

type Cotacao struct {
	Code        string `json:"code"`
	Codein      string `json:"codein"`
	Name        string `json:"name"`
	High        string `json:"high"`
	Low         string `json:"low"`
	VarBid      string `json:"varBid"`
	PctChange   string `json:"pctChange"`
	Bid         string `json:"bid"`
	Ask         string `json:"ask"`
	Timestamp   string `json:"timestamp"`
	Create_date string `json:"create_date"`
}

type Dolar struct {
	Valor float64
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

		val, err := buscarDadosApi()

		if err != nil {
			log.Println(err)
			panic(err)
		}

		err = salvar(val.Valor)
		if err != nil {
			panic(err)
		}
		w.Write([]byte("Process ok"))
	case <-ctx.Done():
		log.Println("Request cancelada pele cliente")
	}

}

func buscarDadosApi() (*Dolar, error) {

	url := "https://economia.awesomeapi.com.br/json/last/USD-BRL"

	resp, err := http.Get(url)
	if err != nil {
		log.Fatalf("Erro ao fazer a requisição: %v", err)
		return nil, err
	}
	defer resp.Body.Close()

	res, err := io.ReadAll(io.Reader(resp.Body))
	if err != nil {
		log.Fatalf("Erro ao ler o corpo da resposta: %v", err)
		return nil, err
	}

	var dados Dados
	if err := json.Unmarshal(res, &dados); err != nil {
		log.Fatalf("Erro ao fazer parsing do JSON: %v", err)
		return nil, err
	}

	dolar, err := strconv.ParseFloat(dados.USDBRL.Bid, 64)
	if err != nil {
		panic(err)
	}
	return &Dolar{Valor: dolar}, nil
}

func salvar(valor float64) error {
	_, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()
	db, err := sql.Open("sqlite3", "./cotacao.db")
	if err != nil {
		log.Fatalf("Erro ao abrir o banco de dados: %v", err)
		return err
	}
	defer db.Close()

	//Criar tabela cotacao
	sqlStmt := `
	CREATE TABLE IF NOT EXISTS cotacao (
		id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		valor FLOAT
	);
	`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		log.Fatalf("Erro ao criar tabela: %v", err)
		return err
	}

	// Inserir dados em cotacao
	_, err = db.Exec("INSERT INTO cotacao (valor) VALUES (?)", valor)
	if err != nil {
		log.Fatalf("Erro ao inserir dados: %v", err)
	}
	return nil
}
