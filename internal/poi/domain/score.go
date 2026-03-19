package domain

func ScorePOI(tags map[string]string) int {
	score := 0

	if tags == nil {
		return score
	}

	if has(tags, "name") {
		score += 20
	}
	if has(tags, "wikipedia") {
		score += 40
	}
	if has(tags, "wikidata") {
		score += 30
	}
	if hasAny(tags, "website", "contact:website") {
		score += 10
	}
	if has(tags, "opening_hours") {
		score += 5
	}
	if hasAny(tags, "phone", "contact:phone") {
		score += 3
	}
	if has(tags, "image") {
		score += 8
	}
	if has(tags, "heritage") {
		score += 15
	}

	switch {
	case tags["tourism"] == "attraction":
		score += 30
	case tags["tourism"] == "museum":
		score += 28
	case tags["tourism"] != "":
		score += 20
	case tags["historic"] != "":
		score += 20
	case tags["amenity"] == "museum":
		score += 25
	case tags["amenity"] == "library":
		score += 10
	case tags["amenity"] == "theatre":
		score += 15
	case tags["amenity"] == "restaurant":
		score += 8
	case tags["amenity"] == "cafe":
		score += 6
	case tags["leisure"] != "":
		score += 8
	case tags["shop"] != "":
		score += 4
	}

	if n := len(tags); n > 0 {
		if n > 15 {
			n = 15
		}
		score += n
	}

	if !has(tags, "name") {
		score -= 10
	}

	return score
}

func has(tags map[string]string, key string) bool {
	v, ok := tags[key]
	return ok && v != ""
}

func hasAny(tags map[string]string, keys ...string) bool {
	for _, k := range keys {
		if has(tags, k) {
			return true
		}
	}
	return false
}
