package gowave

func (w *Wave) Recenter() error {
	rq := GetRecenterQuery()
	err := w.SendQuery(rq)
	return err
}

func GetRecenterQuery() Query {
	q := Query{
		ID:          "Recenter",
		Type:        "Request",
		PayloadSize: 0,
	}

	return q
}
