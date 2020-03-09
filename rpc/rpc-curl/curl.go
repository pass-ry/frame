package curl

import (
	"bytes"
	"context"
	"net/http"
	"time"

	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"gitlab.ifchange.com/data/cordwood/log"
)

var (
	curlRequestCount = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "curl_request_count",
			Help: "curl request count",
		},
		[]string{"url", "c", "m"},
	)

	curlRequestDuration = promauto.NewSummaryVec(
		prometheus.SummaryOpts{
			Name: "curl_request_duration",
			Help: "curl request duration",
		},
		[]string{"url", "c", "m"},
	)

	curlRequestFailCount = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "curl_request_fail_count",
			Help: "curl request fail count",
		},
		[]string{"url", "c", "m"},
	)
)

func Curl(url string, req *Request, rsp *Response) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(5)*time.Minute)
	defer cancel()
	return CurlWithContext(ctx,
		url, req, rsp)
}

func CurlWithContext(ctx context.Context,
	url string, req *Request, rsp *Response) (err error) {
	start := time.Now()
	c := req.GetC()
	m := req.GetM()
	var sourceErr error
	defer func() {
		curlRequestCount.WithLabelValues(url, c, m).Inc()
		elapsed := time.Since(start)
		curlRequestDuration.WithLabelValues(url, c, m).Observe((float64)(elapsed / time.Millisecond))
		if err != nil && sourceErr != nil {
			curlRequestFailCount.WithLabelValues(url, c, m).Inc()
			log.Error(err)
			return
		}
		log.Debugf("Call %s ErrNo %d ErrMsg %s Request %v Response %v Cost:%v",
			url, rsp.GetErrNo(), rsp.GetErrMsg(), req, rsp, elapsed)
		if rsp.GetErrNo() != 0 {
			curlRequestFailCount.WithLabelValues(url, c, m).Inc()
		}
	}()
	if req == nil || rsp == nil {
		sourceErr = errors.Errorf("Call %s NIL Req Or Rsp", url)
		return sourceErr
	}
	body := bytes.NewBuffer(nil)
	_, err = req.WriteTo(body)
	if err != nil {
		sourceErr = err
		return errors.Wrap(err, "Payload Request")
	}
	var (
		httpRequest  *http.Request
		httpResponse *http.Response
	)

	httpRequest, err = http.NewRequestWithContext(ctx,
		http.MethodPost, url, body)
	if err != nil {
		err = errors.Wrap(err, "http.NewRequestWithContext")
		sourceErr = err
		return
	}
	httpRequest.Header.Add("Content-Type", "application/json")
	httpResponse, err = new(http.Client).Do(httpRequest)
	if err != nil {
		sourceErr = err
		err = errors.Errorf("Call %s %v Request %v", url, err, req)
		return
	}
	defer httpResponse.Body.Close()
	if httpResponse.StatusCode != http.StatusOK {
		sourceErr = errors.Errorf("HTTPStatus:%s", httpResponse.Status)
		err = errors.Errorf("Call %s HTTPStatus:%s Request %v", url, httpResponse.Status, req)
		return
	}
	_, err = rsp.ReadFrom(httpResponse.Body)
	if err != nil {
		sourceErr = err
		err = errors.Errorf("Call %s %v Request %v Response %v",
			url, err, req, rsp)
		return
	}
	return
}
