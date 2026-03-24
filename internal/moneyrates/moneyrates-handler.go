package moneyrates

import (
	"encoding/json"
	"net/http"
)

type MoneyRatesHandler struct {
	service Service
}

// Créer un nouveau "Handler" qui va contenir les méthodes appelées pour les endpoints concernés
func NewMoneyRatesHandler(service Service) *MoneyRatesHandler {
	return &MoneyRatesHandler{
		service: service,
	}
}

func (h *MoneyRatesHandler) ConvertRequestedAmount(w http.ResponseWriter, r *http.Request) {

	var convertedRate ConvertedRate

	err := json.NewDecoder(r.Body).Decode(&convertedRate)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	result, err := h.service.convertRate(convertedRate.Base, convertedRate.To, convertedRate.Amount)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	json.NewEncoder(w).Encode(result)
}
