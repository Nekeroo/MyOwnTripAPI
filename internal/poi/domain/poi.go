package domain

type SearchRequest struct {
	CountryCode string
	City        string
	Limit       int
}

type Place struct {
	OSMType     string
	OSMID       int64
	DisplayName string
	Lat         string
	Lon         string
}

type POI struct {
	ID       int64             `json:"id"`
	OSMType  string            `json:"osm_type"`
	Name     string            `json:"name"`
	Category string            `json:"category"`
	Subtype  string            `json:"subtype"`
	Lat      float64           `json:"lat"`
	Lon      float64           `json:"lon"`
	Address  string            `json:"address,omitempty"`
	Tags     map[string]string `json:"tags,omitempty"`
	Score    int               `json:"score"`
}
