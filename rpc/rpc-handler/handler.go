package handler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"runtime"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"gitlab.ifchange.com/data/cordwood/log"
)

var _ error = (*HandlerError)(nil)

type (
	Handler func(req *Request, rsp *Response) error

	HandlerError struct {
		no  int
		msg string
		err error
	}
)

func WrapError(err error, errNo int, errMsg string) error {
	return &HandlerError{
		no:  errNo,
		msg: errMsg,
		err: err,
	}
}

func (req *Request) Unmarshal(params interface{}) error {
	err := json.Unmarshal(req.GetP(), params)
	if err != nil {
		return WrapError(fmt.Errorf("Unmarshal Json Error %v %+v %s",
			err, params, string(req.GetP())),
			1, "未知参数")
	}
	return nil
}

var (
	rpcRequestCount = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "rpc_request_count",
			Help: "rpc request count",
		},
		[]string{"endpoint"},
	)
	rpcRequestDuration = promauto.NewSummaryVec(
		prometheus.SummaryOpts{
			Name: "rpc_request_duration",
			Help: "rpc request duration",
		},
		[]string{"endpoint"},
	)
	rpcRequestFailCount = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "rpc_request_fail_count",
			Help: "rpc request fail count",
		},
		[]string{"endpoint", "err"},
	)
)

func Wrap(handler Handler) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// safe body close
		defer func() {
			if r.Body != nil {
				r.Body.Close()
			}
		}()

		// add response header
		w.Header().Add("Content-Type", "application/json")

		var (
			req = NewRequest()
			rsp = NewResponse()
			err error
		)

		// metrics
		start := time.Now()
		path := r.URL.Path
		rpcRequestCount.WithLabelValues(path).Inc()
		defer func() {
			elapsed := (float64)(time.Since(start) / time.Millisecond)
			rpcRequestDuration.WithLabelValues(path).Observe(elapsed)
			log.Infof("HTTP URI %s Cost %v Response %v",
				r.RequestURI, time.Since(start), rsp)
			if rsp.GetErrNo() != 0 {
				msg := fmt.Sprintf("err_no:%d err_msg:%s",
					rsp.GetErrNo(), rsp.GetErrMsg())
				rpcRequestFailCount.WithLabelValues(path, msg).Inc()
			}
		}()

		// rpc response
		defer func() {
			// force use same header
			rsp.Header = req.Header
			// write to response body
			w.WriteHeader(http.StatusOK)
			_, err := rsp.WriteTo(w)
			if err == nil {
				return
			}
			log.Errorf("HTTP Handle Error URI %s Write Response Error %v Request %v Response %v",
				r.RequestURI, err, req, rsp)
			w.WriteHeader(http.StatusServiceUnavailable)
		}()

		// wrap rpc error
		defer func() {
			if err == nil {
				return
			}

			var (
				errNo  int    = -1
				errMsg string = "系统错误"
			)

			if handlerError, ok := err.(*HandlerError); ok {
				errNo = handlerError.no
				errMsg = handlerError.msg
			}
			rsp.SetErrNo(errNo)
			rsp.SetErrMsg(errMsg)

			log.Errorf("HTTP Handle Error URI %s Error %v Request %v Response %v",
				r.RequestURI, err, req, rsp)
		}()

		// http method protect
		if r.Method != http.MethodPost {
			err = WrapError(fmt.Errorf("Allow POST Only"), 1, "参数错误")
			return
		}

		// rpc request limit
		body, readErr := ioutil.ReadAll(r.Body)
		if readErr != nil {
			err = WrapError(fmt.Errorf("Read Body Error %v", readErr),
				1, "参数错误")
			return
		}

		log.Infof("HTTP URI %s Request %s",
			r.RequestURI, body)

		// rpc request unmarshal
		unmarshalErr := json.Unmarshal(body, req)
		if unmarshalErr != nil {
			err = WrapError(fmt.Errorf("Unmarshal Error %v", unmarshalErr),
				1, "参数错误")
			return
		}

		// handler recover
		defer func() {
			recoverErr := recover()
			if recoverErr == nil {
				return
			}
			buf := make([]byte, 1<<10)
			num := runtime.Stack(buf, false)
			err = WrapError(fmt.Errorf("Router %s PANIC %v %v %v",
				r.RequestURI, recoverErr, num, string(buf)),
				-1, "系统错误")
		}()

		// handler exec
		err = handler(req, rsp)
	}
}

func (err *HandlerError) Error() string {
	return fmt.Sprintf("Response Error No:%d Msg:%s Err:%v",
		err.no, err.msg, err.err)
}
