package clients

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"nekeroo/myowntrip/internal/poi/domain"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Geocoder interface {
	SearchPlace(ctx context.Context, city, countryCode string) (*domain.Place, error)
}

type NominatimClient struct {
	BaseURL    string
	HTTPClient *http.Client
	UserAgent  string
}

type nominatimPlace struct {
	OSMType     string `json:"osm_type"`
	OSMID       int64  `json:"osm_id"`
	DisplayName string `json:"display_name"`
	Lat         string `json:"lat"`
	Lon         string `json:"lon"`
}

func NewNominatimClient(userAgent string) *NominatimClient {
	return &NominatimClient{
		BaseURL:    "https://nominatim.openstreetmap.org",
		HTTPClient: &http.Client{Timeout: 20 * time.Second},
		UserAgent:  userAgent,
	}
}

func (c *NominatimClient) SearchPlace(ctx context.Context, city, countryCode string) (*domain.Place, error) {
	params := url.Values{}
	params.Set("format", "jsonv2")
	params.Set("limit", "1")
	params.Set("email", "m.grattardpro@gmail.com")

	if city != "" {
		params.Set("city", city)
	}
	if countryCode != "" {
		params.Set("countrycodes", strings.ToLower(countryCode))
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		c.BaseURL+"/search?"+params.Encode(),
		nil,
	)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "my-poi-api/1.0 (contact: dev@tondomaine.com)")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Accept-Language", "fr")

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(io.LimitReader(res.Body, 4096))
		return nil, fmt.Errorf("nominatim returned status=%d body=%s", res.StatusCode, string(body))
	}

	var out []nominatimPlace
	if err := json.NewDecoder(res.Body).Decode(&out); err != nil {
		return nil, err
	}
	if len(out) == 0 {
		return nil, fmt.Errorf("no place found")
	}

	return &domain.Place{
		OSMType:     out[0].OSMType,
		OSMID:       out[0].OSMID,
		DisplayName: out[0].DisplayName,
		Lat:         out[0].Lat,
		Lon:         out[0].Lon,
	}, nil
}
