package main

import (
	gowave "github.com/AMcPherran/go-wave"
)

func decodeByteStream(data []byte) gowave.Query {
	q := gowave.NewQuery(data)
	return q
}
