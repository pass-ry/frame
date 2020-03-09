package handler

import (
	"encoding/json"
	"io"

	"github.com/pkg/errors"
	header "gitlab.ifchange.com/data/cordwood/rpc/rpc-header"
)

var (
	_ io.WriterTo = (*Response)(nil)
)

func NewResponse() *Response {
	return &Response{Header: header.NewEmptyHeader()}
}

type Response struct {
	header.Header `json:"header"`
	Response      struct {
		ErrNo   int         `json:"err_no"`
		ErrMsg  string      `json:"err_msg"`
		Results interface{} `json:"results"`
	} `json:"response"`
}

func (rsp *Response) String() string {
	output, _ := json.Marshal(rsp)
	if len(output) > 1000 {
		return string(append(output[:1000], []byte("......")...))
	}
	return string(output)
}

func (rsp *Response) GetErrNo() int                  { return rsp.Response.ErrNo }
func (rsp *Response) SetErrNo(errNo int)             { rsp.Response.ErrNo = errNo }
func (rsp *Response) GetErrMsg() string              { return rsp.Response.ErrMsg }
func (rsp *Response) SetErrMsg(errMsg string)        { rsp.Response.ErrMsg = errMsg }
func (rsp *Response) GetResults() interface{}        { return rsp.Response.Results }
func (rsp *Response) SetResults(results interface{}) { rsp.Response.Results = results }

func (rsp *Response) WriteTo(body io.Writer) (int64, error) {
	rspJson, err := json.Marshal(rsp)
	if err != nil {
		return 0, errors.Wrap(err, "json.Marshal")
	}
	size, err := body.Write(rspJson)
	if err != nil {
		return 0, errors.Wrap(err, "Write body")
	}
	return int64(size), nil
}
