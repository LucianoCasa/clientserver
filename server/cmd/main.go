package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/lucianocasa/clientserver/server/internal/db"
	"github.com/lucianocasa/clientserver/server/internal/repository"
	"github.com/lucianocasa/clientserver/server/internal/service"
	"github.com/lucianocasa/clientserver/server/pkg/exchange"

	"github.com/joho/godotenv"
)

type cotacaousdbrl struct {
	usdbrl float64
}

func BuscaCotacao(w http.ResponseWriter, r *http.Request) {
	if err := godotenv.Load(); err != nil {
		log.Println("Erro carregando o arquivo .env")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	timeoutReq, err := strconv.Atoi(os.Getenv("API_REQUEST_QUOTE_TIMEOUT_MS"))
	if err != nil {
		log.Println("Erro na busca do timeout do request.", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	timeoutBd, err := strconv.Atoi(os.Getenv("API_DB_QUOTE_TIMEOUT_MS"))
	if err != nil {
		log.Println("Erro na busca do timeout do BD.", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	clientCtx := r.Context() // Se o client cancelar a requisição para
	ctxApi, cancel := context.WithTimeout(clientCtx, time.Duration(timeoutReq)*time.Millisecond)
	defer cancel()

	log.Println("Solicitação de cotação")

	ret, err := exchange.CotacaoUSDBRL(ctxApi)
	if err != nil {
		// Verifica se o erro foi causado pelo cancelamento do cliente
		if clientCtx.Err() == context.Canceled {
			log.Println("Requisição cancelada pelo cliente")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// Verifica se o erro foi causado pelo timeout do servidor
		if ctxApi.Err() == context.DeadlineExceeded {
			log.Println("Erro: o tempo para buscar a cotação foi excedido (timeout do servidor).")
			w.WriteHeader(http.StatusGatewayTimeout)
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Erro ao buscar cotação:", err)
		return
	}

	var parsed map[string]float64
	err = json.Unmarshal(ret, &parsed)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Erro ao parsear JSON da cotação:", err)
		return
	}

	valor, ok := parsed["usdbrl"]
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Campo 'usdbrl' não encontrado na resposta")
		return
	}

	log.Printf("Cotação USD/BRL: %.4f", valor)

	// Salvar no banco
	err = SalvaBDCotacao(clientCtx, timeoutBd, valor)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Erro ao salvar a cotação:", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(ret)
	log.Println("Requisição finalizada")
}

func SalvaBDCotacao(ctx context.Context, timeoutBd int, valor float64) error {
	ctx, cancel := context.WithTimeout(ctx, time.Duration(timeoutBd)*time.Millisecond)
	defer cancel()

	db.Init("cotacao.db")
	repo := repository.NewCotacaoRepository(db.DBConn)
	svc := service.NewCotacaoService(repo)

	err := svc.CriarCotacao(ctx, valor)
	if err != nil {
		return err
	}
	return nil
}

func ListaBDCotacao() {
	db.Init("cotacao.db")
	repo := repository.NewCotacaoRepository(db.DBConn)
	svc := service.NewCotacaoService(repo)

	cotacoes, err := svc.ListarCotacoes()
	if err != nil {
		log.Fatalf("Erro ao listar cotação: %v", err)
	}
	log.Println("Cotações listadas: ", cotacoes)
}

func start() {
	log.Println("Servidor no ar.")
	http.HandleFunc("/", BuscaCotacao)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalf("Erro ao iniciar o servidor: %v", err)
	}
}

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("Uso: go run main.go [start|list|migrate up|migrate down]")
	}

	switch os.Args[1] {
	case "start":
		start()
	case "list":
		ListaBDCotacao()
	case "migrate":
		if len(os.Args) < 3 {
			log.Fatalf("Uso: go run main.go migrate [up|down]")
		}
		switch os.Args[2] {
		case "up":
			db.Init("cotacao.db")
			err := db.MigrateUp(db.DBConn)
			if err != nil {
				log.Fatalf("Erro ao executar migrate up: %v", err)
			}
			log.Println("Migração executada com sucesso (up)")
		case "down":
			db.Init("cotacao.db")
			err := db.MigrateDown(db.DBConn)
			if err != nil {
				log.Fatalf("Erro ao executar migrate down: %v", err)
			}
			log.Println("Migração revertida com sucesso (down)")
		default:
			fmt.Println("Comando inválido para migrate. Use 'up' ou 'down'.")
		}
	default:
		fmt.Println("Comando desconhecido:", os.Args[1])
	}
}
