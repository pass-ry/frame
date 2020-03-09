package json

import (
	"bytes"
	"encoding/json"
)

type (
	EncodeOpt uint8
	DecodeOpt uint8
)

const (
	unEscapeHTML EncodeOpt = 1 << iota
)

func Marshal(v interface{}, opts ...EncodeOpt) ([]byte, error) {
	if len(opts) == 0 {
		return json.Marshal(v)
	}
	var flag EncodeOpt
	for _, opt := range opts {
		flag |= opt
	}
	bytes := bytes.NewBuffer([]byte{})
	jsonEncoder := json.NewEncoder(bytes)
	if flag&unEscapeHTML != 0 {
		jsonEncoder.SetEscapeHTML(false)
	}

	err := jsonEncoder.Encode(v)
	return bytes.Bytes(), err
}

func UnEscapeHTML() EncodeOpt {
	return unEscapeHTML
}

func Unmarshal(data []byte, v interface{}, opts ...DecodeOpt) error {
	return json.Unmarshal(data, v)
}
