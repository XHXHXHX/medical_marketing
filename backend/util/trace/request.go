package trace

import (
	"net/http"
)

func GetTraceFromHTTP(req *http.Request) string {
	if req == nil {
		return ""
	}
	return req.Header.Get(IotTraceHeader)
}

func SetHTTPNewTraceID(req *http.Request) {
	if req == nil {
		return
	}

	req.Header.Set(IotTraceHeader, NewTraceID())
}

func SetHTTPTraceID(req *http.Request, traceID string) {
	if req == nil {
		return
	}

	req.Header.Set(IotTraceHeader, traceID)
}

func SetHTTPResponseTraceID(w http.ResponseWriter, traceID string) {
	w.Header().Set(IotTraceHeader, traceID)
}
