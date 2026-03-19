package moneyrates

import (
	"log"
	"nekeroo/myowntrip/internal/json"
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

// Méthode appelée lorsqu'on veut récupérer le cours actuel de l'euro à l'étranger
func (h *MoneyRatesHandler) RetrieveLatestMoneyRates(w http.ResponseWriter, r *http.Request) {
	moneyRates, err := h.service.retrieveLatestRates()

	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.Write(w, http.StatusOK, moneyRates)
}
