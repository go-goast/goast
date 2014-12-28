package main

import (
	"encoding/json"
)

func (t *User) Json() string {
	b, err := json.Marshal(t)
	if err != nil {
		return err.Error()
	}
	return string(b)
}
