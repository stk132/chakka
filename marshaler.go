package main

import (
	"encoding/json"
	"github.com/fireworq/fireworq/model"
)

type marshaler interface {
	marshal() ([]byte, error)
}

type myQueue model.Queue
type myRouting model.Routing

func (m myQueue) marshal() ([]byte, error) {
	return json.Marshal(&m)
}

func (m myRouting) marshal() ([]byte, error) {
	return json.Marshal(&m)
}
