package httppoi

import (
	"encoding/json"
	"nekeroo/myowntrip/internal/poi/domain"
	"nekeroo/myowntrip/internal/poi/service"
	"net/http"
	"strconv"
)

type Handler struct {
	Service *service.Service
}

func NewHandler(s *service.Service) *Handler {
	return &Handler{Service: s}
}

func (h *Handler) SearchPOIs(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	query := r.URL.Query()
	city := query.Get("city")
	countryCode := query.Get("countryCode")

	limit := 50
	if raw := query.Get("limit"); raw != "" {
		if parsed, err := strconv.Atoi(raw); err == nil {
			limit = parsed
		}
	}

	req := domain.SearchRequest{
		City:        city,
		CountryCode: countryCode,
		Limit:       limit,
	}

	pois, err := h.Service.Search(ctx, req)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"count": len(pois),
		"items": pois,
	})
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}
