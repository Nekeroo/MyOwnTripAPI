package moneyrates

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type MoneyRates struct {
	Amount float64            `json:"amount"`
	Base   string             `json:"base"`
	Date   string             `json:"date"`
	Rates  map[string]float64 `json:"rates"`
}

type Service struct {
	urlFrankfurter string
}

func NewService() *Service {
	return &Service{
		urlFrankfurter: "http://localhost:8082/v1",
	}
}

func (s *Service) retrieveLatestRates() (MoneyRates, error) {
	res, err := http.Get(s.urlFrankfurter + "/latest")
	if err != nil {
		return MoneyRates{}, fmt.Errorf("error calling Frankfurter API: %w", err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return MoneyRates{}, fmt.Errorf("failed to read response body: %w", err)
	}

	log.Println(string(body))

	var rates MoneyRates
	if err := json.Unmarshal(body, &rates); err != nil {
		return MoneyRates{}, fmt.Errorf("failed to parse JSON: %w", err)
	}

	return rates, nil
}
