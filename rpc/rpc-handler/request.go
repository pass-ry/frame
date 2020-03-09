package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"

	"github.com/pkg/errors"
	header "gitlab.ifchange.com/data/cordwood/rpc/rpc-header"
)

var (
	_ io.ReaderFrom = (*Request)(nil)
)

func NewRequest() *Request {
	return &Request{Header: header.NewEmptyHeader()}
}

type Request struct {
	header.Header `json:"header"`
	Request       struct {
		C string          `json:"c"`
		M string          `json:"m"`
		P json.RawMessage `json:"p"`
	} `json:"request"`
}

func (req *Request) String() string {
	output, _ := json.Marshal(req)
	if len(output) > 1000 {
		return string(append(output[:1000], []byte("......")...))
	}
	return string(output)
}

func (req *Request) GetC() string           { return req.Request.C }
func (req *Request) SetC(c string)          { req.Request.C = c }
func (req *Request) GetM() string           { return req.Request.M }
func (req *Request) SetM(m string)          { req.Request.M = m }
func (req *Request) GetP() json.RawMessage  { return req.Request.P }
func (req *Request) SetP(p json.RawMessage) { req.Request.P = p }

func (req *Request) ReadFrom(body io.Reader) (int64, error) {
	if req == nil {
		return 0, fmt.Errorf("Nil payload.Request")
	}
	bodyBytes, err := ioutil.ReadAll(body)
	if err != nil {
		return 0, errors.Wrap(err, "read body")
	}
	err = json.Unmarshal(bodyBytes, &req)
	if err != nil {
		return 0, errors.Wrap(err, "json.Unmarshal")
	}
	return int64(len(bodyBytes)), nil
}
