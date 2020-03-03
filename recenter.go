package gowave

func GetRecenterQuery() Query {
	q := Query{
		ID:          "Recenter",
		Type:        "Request",
		PayloadSize: 0,
	}

	return q
}
