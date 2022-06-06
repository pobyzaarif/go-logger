package http

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httputil"
	"strings"
)

func DumpRequest(req *http.Request, hiddenHeaders []string) string {
	if req == nil {
		return ""
	}

	req.Clone(context.TODO())

	requestDump, err := httputil.DumpRequest(req, true)
	if err != nil {
		return fmt.Sprintf("%+v", req)
	}

	stringDump := string(requestDump)

	for _, hiddenHeader := range hiddenHeaders {
		val := req.Header.Get(hiddenHeader)
		if val != "" {
			stringDump = strings.Replace(stringDump, val, "**hidden**", -1)
		}
	}

	return stringDump
}

func DumpResponse(resp *http.Response) string {
	// Handling nil pointer
	if resp == nil {
		return ""
	}

	responseDump, err := httputil.DumpResponse(resp, true)
	if err != nil {
		return fmt.Sprintf("%+v", resp)
	}
	return string(responseDump)
}

// NetworkLog : network log wrapper
func NetworkLog(data map[string]interface{}) map[string]interface{} {
	return map[string]interface{}{
		"__gologger__": 1,
		"net":          data,
	}
}
