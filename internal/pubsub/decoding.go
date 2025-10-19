package pubsub

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
)

func decodeJSON[T any](b []byte) (T, error) {
	var v T
	err := json.Unmarshal(b, &v)
	return v, err
}

func decodeGob[T any](b []byte) (T, error) {
	var v T
	buff := bytes.NewBuffer(b)
	dec := gob.NewDecoder(buff)
	err := dec.Decode(&v)
	return v, err
}
