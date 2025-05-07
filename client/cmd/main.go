package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/lucianocasa/clientserver/client/pkg/utils/fileutils"

	"github.com/joho/godotenv"
)

type Cotacao struct {
	USDBRL float64 `json:"usdbrl"`
}

func main() {
	var cotacao Cotacao

	if err := godotenv.Load(); err != nil {
		log.Println("Erro carregando o arquivo .env")
		return
	}

	var arquivoCotacao string = os.Getenv("QUOTE_FILE_LOCATION")
	var urlCotacao string = os.Getenv("API_URL")

	timeoutReq, err := strconv.Atoi(os.Getenv("API_REQUEST_TIMEOUT_MS"))
	if err != nil {
		log.Println("Erro na busca do timeout.", err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeoutReq)*time.Millisecond)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", urlCotacao, nil)
	if err != nil {
		log.Println("Erro na requisição: Acesso ao Servidor não funcionou.", err)
		return
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			log.Println("Erro: o tempo para buscar a cotação foi excedido (timeout).")
		} else {
			log.Println("Erro ao enviar requisição:", err)

		}
		return
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Println("Erro na requisição: Erro ao ler o retorno.")
		return
	}

	err = json.Unmarshal(body, &cotacao)
	if err != nil {
		log.Println("Erro ao fazer parse:", err)
		return
	}

	fileutils.WriteNewFile(arquivoCotacao, "Dólar: "+fmt.Sprintf("%.4f", cotacao.USDBRL))
	log.Println("Arquivo (" + arquivoCotacao + ") gravado com sucesso!")

}
