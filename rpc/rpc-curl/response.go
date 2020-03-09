package curl

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"

	"github.com/pkg/errors"
	header "gitlab.ifchange.com/data/cordwood/rpc/rpc-header"
)

var (
	_ io.ReaderFrom = (*Response)(nil)
)

func NewResponse() *Response {
	return &Response{Header: header.NewEmptyHeader()}
}

type Response struct {
	header.Header `json:"header,omitempty"`
	Response      struct {
		ErrNo   int             `json:"err_no"`
		ErrMsg  string          `json:"err_msg"`
		Results json.RawMessage `json:"results"`
	} `json:"response"`
}

func (rsp *Response) String() string {
	output, _ := json.Marshal(rsp)
	if len(output) > 1000 {
		return string(append(output[:1000], []byte("......")...))
	}
	return string(output)
}

func (rsp *Response) GetErrNo() int                      { return rsp.Response.ErrNo }
func (rsp *Response) SetErrNo(errNo int)                 { rsp.Response.ErrNo = errNo }
func (rsp *Response) GetErrMsg() string                  { return rsp.Response.ErrMsg }
func (rsp *Response) SetErrMsg(errMsg string)            { rsp.Response.ErrMsg = errMsg }
func (rsp *Response) GetResults() json.RawMessage        { return rsp.Response.Results }
func (rsp *Response) SetResults(results json.RawMessage) { rsp.Response.Results = results }

func (rsp *Response) ReadFrom(body io.Reader) (int64, error) {
	if rsp == nil {
		return 0, fmt.Errorf("Nil payload.Response")
	}
	bodyBytes, err := ioutil.ReadAll(body)
	if err != nil {
		return 0, errors.Wrap(err, "read body")
	}
	err = json.Unmarshal(bodyBytes, &rsp)
	if err != nil {
		return 0, errors.Wrap(err, "json.Unmarshal")
	}
	return int64(len(bodyBytes)), nil
}
