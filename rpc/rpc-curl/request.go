package curl

import (
	"encoding/json"
	"io"

	"github.com/pkg/errors"
	header "gitlab.ifchange.com/data/cordwood/rpc/rpc-header"
)

var (
	_ io.WriterTo = (*Request)(nil)
)

func NewRequest() *Request {
	return &Request{Header: header.NewHeader()}
}

type Request struct {
	header.Header `json:"header,omitempty"`
	Request       struct {
		C string      `json:"c"`
		M string      `json:"m"`
		P interface{} `json:"p"`
	} `json:"request"`
}

func (req *Request) String() string {
	output, _ := json.Marshal(req)
	if len(output) > 1000 {
		return string(append(output[:1000], []byte("......")...))
	}
	return string(output)
}

func (req *Request) GetC() string       { return req.Request.C }
func (req *Request) SetC(c string)      { req.Request.C = c }
func (req *Request) GetM() string       { return req.Request.M }
func (req *Request) SetM(m string)      { req.Request.M = m }
func (req *Request) GetP() interface{}  { return req.Request.P }
func (req *Request) SetP(p interface{}) { req.Request.P = p }

func (req *Request) WriteTo(body io.Writer) (int64, error) {
	reqJson, err := json.Marshal(req)
	if err != nil {
		return 0, errors.Wrap(err, "json.Marshal")
	}
	size, err := body.Write(reqJson)
	if err != nil {
		return 0, errors.Wrap(err, "Write body")
	}
	return int64(size), nil
}
