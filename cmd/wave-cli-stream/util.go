package main

import (
	gowave "github.com/AMcPherran/go-wave"
)

func decodeByteStream(data []byte) gowave.Query {
	q := gowave.Query{
		ID:          data[0],
		PayloadSize: len(data),
	}

	return q
}
