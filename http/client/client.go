package http

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	xmlToJson "github.com/basgys/goxml2json"
	goLoggerHttp "github.com/pobyzaarif/go-logger/http"
	goLogger "github.com/pobyzaarif/go-logger/logger"
)

type (
	ResponseBodyFormat int

	ProxyConfig struct {
		Host string
		Port int
	}
)

//ResponseBodyFormat possible values
const (
	RawResponseBodyFormat ResponseBodyFormat = iota
	JSONResponseBodyFormat
	XMLResponseBodyFormat
)

var logger = goLogger.NewLog("OUTBOUND_REQUEST")

func Call(
	ctx context.Context,
	request *http.Request,
	timeout time.Duration,
	responseBodyFormat ResponseBodyFormat,
	responseBody interface{},
	proxyConfig *ProxyConfig) (int, error) {
	ctxTrackerID := ctx.Value("tracker_id")
	trackerID := ""
	if ctxTrackerID != nil {
		trackerID = fmt.Sprintf("%v", ctxTrackerID)
	}
	logger.SetTrackerID(trackerID)

	var client http.Client
	client.Timeout = timeout

	httpLog := map[string]interface{}{
		"host":               request.URL.Host,
		"method":             request.Method,
		"url":                request.URL.Path,
		"request":            goLoggerHttp.DumpRequest(request, []string{"Authorization"}),
		"response":           "",
		"response_http_code": 0,
	}

	if proxyConfig != nil {
		proxyUrl, err := url.Parse(fmt.Sprintf("http://%s:%d", proxyConfig.Host, proxyConfig.Port))
		if err != nil {
			logger.ErrorWithData("failed to parse proxy url", goLoggerHttp.NetworkLog(httpLog), err)

			return -1, errors.New("failed to parse proxy url")
		}

		client.Transport = &http.Transport{Proxy: http.ProxyURL(proxyUrl)}
	}

	logger.SetTimerStart(time.Now())

	res, err := client.Do(request)

	httpLog["response"] = goLoggerHttp.DumpResponse(res)

	if err != nil {
		errMessage := "error is " + err.Error()
		urlErr, ok := err.(*url.Error)
		if ok && urlErr.Timeout() {
			logger.ErrorWithData("timeout on request", goLoggerHttp.NetworkLog(httpLog), urlErr)

			return 0, errors.New(errMessage)
		}

		logger.ErrorWithData("failed on request", goLoggerHttp.NetworkLog(httpLog), urlErr)

		return 0, errors.New(errMessage)
	}

	defer res.Body.Close()

	httpLog["response_http_code"] = res.StatusCode
	httpLog = goLoggerHttp.NetworkLog(httpLog) // wrapping logger to standart network log

	buffer, err := ioutil.ReadAll(res.Body)
	if err != nil {
		logger.ErrorWithData("failed to get response body", httpLog, err)

		return res.StatusCode, errors.New("failed to get response body")
	}

	if responseBodyFormat == XMLResponseBodyFormat {
		jsonPresenter, err := xmlToJson.Convert(bytes.NewReader(buffer))
		if err != nil {
			logger.ErrorWithData("failed to convert xml to json", httpLog, err)

			return res.StatusCode, errors.New("failed to convert xml to json")
		}

		err = json.Unmarshal(jsonPresenter.Bytes(), responseBody)
		if err != nil {
			logger.ErrorWithData("failed to parsing json from xml response", httpLog, err)

			return res.StatusCode, errors.New("failed to parsing json from xml response")
		}
	} else if responseBodyFormat == JSONResponseBodyFormat {
		err = json.Unmarshal(buffer, responseBody)
		if err != nil {
			logger.ErrorWithData("failed to parsing json response", httpLog, err)

			return res.StatusCode, errors.New("failed to parsing json response")
		}
	} else {
		responseBody = buffer
	}

	logger.InfoWithData("success", httpLog)

	return res.StatusCode, nil
}
