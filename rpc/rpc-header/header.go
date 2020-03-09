package header

import (
	"fmt"

	"gitlab.ifchange.com/data/cordwood/util/random"
)

type Header map[string]interface{}

func (h Header) GetLogID() string {
	return fmt.Sprintf("%s", h["log_id"])
}

func (h Header) SetLogID(logID string) {
	h["log_id"] = logID
}

func NewEmptyHeader() Header {
	return map[string]interface{}{
		"log_id": random.RandStr(16),
	}
}

func NewHeader() Header {
	return map[string]interface{}{
		"appid":     22,
		"app_id":    22,
		"log_id":    random.RandStr(16),
		"uid":       "",
		"uname":     "",
		"provider":  "data",
		"signid":    "",
		"version":   "",
		"ip":        "",
		"timestamp": "",
		"nonce":     "",
		"signature": "",
	}
}
