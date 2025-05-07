package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/lucianocasa/clientserver/server/internal/model"
)

type CotacaoRepository struct {
	DBConn *sql.DB
}

func NewCotacaoRepository(db *sql.DB) *CotacaoRepository {
	return &CotacaoRepository{DBConn: db}
}

func (r *CotacaoRepository) Criar(ctx context.Context, c model.Cotacao) error {
	_, err := r.DBConn.ExecContext(ctx, "INSERT INTO cotacoes (valor, created_at) VALUES (?, ?)", c.Valor, c.CreatedAt.Format(time.RFC3339))
	return err
}

func (r *CotacaoRepository) Listar() ([]model.Cotacao, error) {
	rows, err := r.DBConn.Query("SELECT id, valor, created_at FROM cotacoes")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cotacoes []model.Cotacao
	for rows.Next() {
		var c model.Cotacao
		var createdAtStr string
		if err := rows.Scan(&c.ID, &c.Valor, &createdAtStr); err != nil {
			return nil, err
		}
		c.CreatedAt, _ = time.Parse(time.RFC3339, createdAtStr)
		cotacoes = append(cotacoes, c)
	}
	return cotacoes, nil
}
