package clients

import (
	"bytes"
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

type POIFetcher interface {
	SearchPOIsInArea(ctx context.Context, areaName string, maxFetch int) ([]domain.POI, error)
}

type OverpassClient struct {
	BaseURL    string
	HTTPClient *http.Client
	UserAgent  string
}

type overpassResponse struct {
	Elements []overpassElement `json:"elements"`
}

type overpassElement struct {
	Type   string            `json:"type"`
	ID     int64             `json:"id"`
	Lat    float64           `json:"lat,omitempty"`
	Lon    float64           `json:"lon,omitempty"`
	Center *overpassCenter   `json:"center,omitempty"`
	Tags   map[string]string `json:"tags,omitempty"`
}

type overpassCenter struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}

func NewOverpassClient(userAgent string) *OverpassClient {
	return &OverpassClient{
		BaseURL:    "https://overpass-api.de/api/interpreter",
		HTTPClient: &http.Client{Timeout: 45 * time.Second},
		UserAgent:  userAgent,
	}
}

func (c *OverpassClient) SearchPOIsInArea(ctx context.Context, areaName string, maxFetch int) ([]domain.POI, error) {
	query := buildAreaQuery(areaName, maxFetch)

	form := url.Values{}
	form.Set("data", query)

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		c.BaseURL,
		bytes.NewBufferString(form.Encode()),
	)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", c.UserAgent)
	req.Header.Set("Accept", "application/json")

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(io.LimitReader(res.Body, 4096))
		return nil, fmt.Errorf("overpass returned status=%d body=%s", res.StatusCode, string(body))
	}

	var out overpassResponse
	if err := json.NewDecoder(res.Body).Decode(&out); err != nil {
		return nil, err
	}

	pois := make([]domain.POI, 0, len(out.Elements))
	for _, el := range out.Elements {
		lat, lon := coordinates(el)
		category, subtype := detectCategory(el.Tags)

		pois = append(pois, domain.POI{
			ID:       el.ID,
			OSMType:  el.Type,
			Name:     fallbackName(el.Tags),
			Category: category,
			Subtype:  subtype,
			Lat:      lat,
			Lon:      lon,
			Address:  buildAddress(el.Tags),
			Tags:     el.Tags,
		})
	}

	return pois, nil
}

func buildAreaQuery(areaName string, maxFetch int) string {
	escaped := escapeOverpassString(areaName)

	return fmt.Sprintf(`
[out:json][timeout:25];
area["name"="%s"]["boundary"="administrative"]->.a;
(
  nwr(area.a)["tourism"];
  nwr(area.a)["amenity"];
  nwr(area.a)["shop"];
  nwr(area.a)["leisure"];
  nwr(area.a)["historic"];
);
out center tags %d;
`, escaped, maxFetch)
}

func coordinates(el overpassElement) (float64, float64) {
	if el.Center != nil {
		return el.Center.Lat, el.Center.Lon
	}
	return el.Lat, el.Lon
}

func detectCategory(tags map[string]string) (string, string) {
	for _, key := range []string{"tourism", "amenity", "historic", "leisure", "shop"} {
		if v, ok := tags[key]; ok && v != "" {
			return key, v
		}
	}
	return "unknown", "unknown"
}

func fallbackName(tags map[string]string) string {
	if tags["name"] != "" {
		return tags["name"]
	}
	return "(sans nom)"
}

func buildAddress(tags map[string]string) string {
	parts := []string{
		tags["addr:housenumber"],
		tags["addr:street"],
		tags["addr:postcode"],
		tags["addr:city"],
	}

	out := make([]string, 0, len(parts))
	for _, p := range parts {
		if strings.TrimSpace(p) != "" {
			out = append(out, p)
		}
	}

	return strings.Join(out, ", ")
}

func escapeOverpassString(s string) string {
	s = strings.ReplaceAll(s, `\`, `\\`)
	s = strings.ReplaceAll(s, `"`, `\"`)
	return s
}
