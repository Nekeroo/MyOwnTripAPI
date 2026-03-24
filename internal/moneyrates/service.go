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

type ConvertedRate struct {
	Amount float64 `json:"amount"`
	Base   string  `json:"base"`
	To     string  `json:"to"`
}

type Service struct {
	urlFrankfurter string
}

func NewService() *Service {
	return &Service{
		urlFrankfurter: "https://api.frankfurter.dev/v1",
	}
}

func (s *Service) convertRate(from string, to string, amount float64) (ConvertedRate, error) {
	res, err := http.Get(s.urlFrankfurter + "/latest?base=" + from + "&symbols=" + to)

	if err != nil {
		return ConvertedRate{}, fmt.Errorf("Error calling Frankfurter API : %w", err)
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)

	if err != nil {
		return ConvertedRate{}, fmt.Errorf("Failed to read response : %w", err)
	}

	log.Println(string(body))

	var rates MoneyRates
	if err := json.Unmarshal(body, &rates); err != nil {
		return ConvertedRate{}, fmt.Errorf("failed to parse JSON: %w", err)
	}

	result := ConvertedRate{
		(amount * rates.Rates[to]),
		from,
		to,
	}

	return result, nil
}
