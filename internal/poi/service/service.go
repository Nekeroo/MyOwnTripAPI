package service

import (
	"context"
	"errors"
	"nekeroo/myowntrip/internal/poi/clients"
	"nekeroo/myowntrip/internal/poi/domain"
	"sort"
	"strings"
)

const (
	DefaultLimit    = 50
	MaxReturnedPOIs = 50
	MaxFetchedPOIs  = 300
)

type Service struct {
	Geocoder   clients.Geocoder
	POIFetcher clients.POIFetcher
}

func New(geocoder clients.Geocoder, poiFetcher clients.POIFetcher) *Service {
	return &Service{
		Geocoder:   geocoder,
		POIFetcher: poiFetcher,
	}
}

func (s *Service) Search(ctx context.Context, req domain.SearchRequest) ([]domain.POI, error) {
	if err := validateRequest(req); err != nil {
		return nil, err
	}

	limit := req.Limit
	if limit <= 0 {
		limit = DefaultLimit
	}
	if limit > MaxReturnedPOIs {
		limit = MaxReturnedPOIs
	}

	place, err := s.Geocoder.SearchPlace(ctx, req.City, req.CountryCode)
	if err != nil {
		return nil, err
	}

	pois, err := s.POIFetcher.SearchPOIsInPlace(ctx, *place, MaxFetchedPOIs)

	if err != nil {
		return nil, err
	}

	for i := range pois {
		pois[i].Score = domain.ScorePOI(pois[i].Tags)
	}

	sort.SliceStable(pois, func(i, j int) bool {
		if pois[i].Score == pois[j].Score {
			return pois[i].Name < pois[j].Name
		}
		return pois[i].Score > pois[j].Score
	})

	pois = deduplicate(pois)

	if len(pois) > limit {
		pois = pois[:limit]
	}

	return pois, nil
}

func validateRequest(req domain.SearchRequest) error {
	if strings.TrimSpace(req.City) == "" && strings.TrimSpace(req.CountryCode) == "" {
		return errors.New("city or countryCode is required")
	}
	return nil
}

func extractAreaName(city, displayName string) string {
	if strings.TrimSpace(city) != "" {
		return city
	}
	if displayName != "" {
		parts := strings.Split(displayName, ",")
		if len(parts) > 0 {
			return strings.TrimSpace(parts[0])
		}
	}
	return displayName
}

func deduplicate(pois []domain.POI) []domain.POI {
	seen := make(map[string]struct{}, len(pois))
	out := make([]domain.POI, 0, len(pois))

	for _, p := range pois {
		key := p.OSMType + ":" + p.Name + ":" + p.Category + ":" + p.Subtype
		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = struct{}{}
		out = append(out, p)
	}

	return out
}
