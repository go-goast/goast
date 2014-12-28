package main

import (
	"encoding/json"
)

type T struct{}

func (t *T) Json() string {
	b, err := json.Marshal(t)
	if err != nil {
		return err.Error()
	}
	return string(b)
}
