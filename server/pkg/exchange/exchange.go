package exchange

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Usdbrl struct {
	Code       string `json:"code"`
	Codein     string `json:"codein"`
	Name       string `json:"name"`
	High       string `json:"high"`
	Low        string `json:"low"`
	VarBid     string `json:"varBid"`
	PctChange  string `json:"pctChange"`
	Bid        string `json:"bid"`
	Ask        string `json:"ask"`
	Timestamp  string `json:"timestamp"`
	CreateDate string `json:"create_date"`
}

type ConvMoeda struct {
	Usdbrl Usdbrl `json:"USDBRL"`
}

func CotacaoUSDBRL(ctx context.Context) ([]byte, error) {
	var result ConvMoeda

	if err := godotenv.Load(); err != nil {
		return nil, err
	}

	urlapi := os.Getenv("API_URL_QUOTE")

	req, err := http.NewRequestWithContext(ctx, "GET", urlapi, nil)
	if err != nil {
		return nil, err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	bid, err := strconv.ParseFloat(result.Usdbrl.Bid, 64)
	if err != nil {
		return nil, err
	}

	retorno := map[string]float64{"usdbrl": bid}
	return json.Marshal(retorno)

}
