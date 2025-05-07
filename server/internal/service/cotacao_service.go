package service

import (
	"context"
	"time"

	"github.com/lucianocasa/clientserver/server/internal/model"
	"github.com/lucianocasa/clientserver/server/internal/repository"
)

type CotacaoService struct {
	Repo *repository.CotacaoRepository
}

func NewCotacaoService(repo *repository.CotacaoRepository) *CotacaoService {
	return &CotacaoService{Repo: repo}
}

func (s *CotacaoService) CriarCotacao(ctx context.Context, valor float64) error {
	c := model.Cotacao{
		Valor:     valor,
		CreatedAt: time.Now(),
	}
	return s.Repo.Criar(ctx, c)
}

func (s *CotacaoService) ListarCotacoes() ([]model.Cotacao, error) {
	return s.Repo.Listar()
}
