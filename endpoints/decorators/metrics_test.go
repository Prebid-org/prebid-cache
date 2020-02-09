package decorators

import (
	"github.com/julienschmidt/httprouter"
	pbcmetrics "github.com/prebid/prebid-cache/metrics"
	"github.com/prebid/prebid-cache/metrics/metricstest"
	"github.com/rcrowley/go-metrics"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSuccessMetrics(t *testing.T) {
	var handler = func(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
		w.WriteHeader(200)
	}
	entry := doRequest(handler)

	//metricstest.AssertSuccessMetricsExist(t, entry)
	actualRequestDuration, _ := HT1["gets.backends.duration"]
	actualRequestCount, _ := HT2["gets.backends.request.total"]
	assert.Equalf(t, int64(1), actualRequestCount, "Successful backend request been accounted for in the total get backend request count, expected = 1; actual = %d\n", actualRequestCount)
	assert.Greater(t, actualRequestDuration, 0.00, "Successful put request duration should be greater than zero")
}

/*
func TestBadRequestMetrics(t *testing.T) {
	var handler = func(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
		w.WriteHeader(400)
	}
	entry := doRequest(handler)

	if entry.Request.Count() != 1 {
		t.Errorf("The request should have been counted.")
	}
	if entry.Duration.Count() != 0 {
		t.Errorf("The request duration should not have been counted.")
	}
	if entry.BadRequest.Count() != 1 {
		t.Errorf("A Bad request should have been counted.")
	}
	if entry.Errors.Count() != 0 {
		t.Errorf("No Errors should have been counted.")
	}
}

func TestErrorMetrics(t *testing.T) {
	var handler = func(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
		w.WriteHeader(500)
	}
	entry := doRequest(handler)
	metricstest.AssertErrorMetricsExist(t, entry)
}

func TestNoExplicitHeaderMetrics(t *testing.T) {
	var handler = func(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {}
	entry := doRequest(handler)
	metricstest.AssertSuccessMetricsExist(t, entry)
}

func TestWriteBytesMetrics(t *testing.T) {
	var handler = func(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
		w.Write([]byte("Success"))
	}
	entry := doRequest(handler)
	metricstest.AssertSuccessMetricsExist(t, entry)
}
*/
func doRequest(handler func(http.ResponseWriter, *http.Request, httprouter.Params)) *pbcmetrics.MetricsEntry {
	reg := metrics.NewRegistry()
	entry := pbcmetrics.NewMetricsEntry("foo", reg)
	monitoredHandler := MonitorHttp(handler, entry)
	monitoredHandler(httptest.NewRecorder(), nil, nil)
	return entry
}

/*Define Mock metrics        */
var HT1 map[string]float64
var HT2 map[string]int64

func CreateMockMetrics() *metrics.Metrics {
	HT1 = make(map[string]float64, 6)
	HT1["puts.current_url.duration"] = 0.00
	HT1["gets.current_url.duration"] = 0.00
	HT1["puts.backends.request_duration"] = 0.00
	HT1["puts.backends.request_size_bytes"] = 0.00
	HT1["gets.backends.duration"] = 0.00
	HT1["connections.connections_opened"] = 0.00
	HT1["extra_ttl_seconds"] = 0.00

	HT2 = make(map[string]int64, 16)
	HT2["puts.current_url.request.total"] = 0
	HT2["puts.current_url.request.error"] = 0
	HT2["puts.current_url.request.bad_request"] = 0
	HT2["gets.current_url.request.total"] = 0
	HT2["gets.current_url.request.error"] = 0
	HT2["gets.current_url.request.bad_request"] = 0
	HT2["puts.backends.add"] = 0
	HT2["puts.backends.json"] = 0
	HT2["puts.backends.xml"] = 0
	HT2["puts.backends.invalid_format"] = 0
	HT2["puts.backends.defines_ttl"] = 0
	HT2["puts.backends.request.error"] = 0
	HT2["gets.backends.request.total"] = 0
	HT2["gets.backends.request.error"] = 0
	HT2["gets.backends.request.bad_request"] = 0
	HT2["connections.connection_error.accept"] = 0
	HT2["connections.connection_error.close"] = 0

	return &metrics.Metrics{MetricEngines: []metrics.CacheMetrics{&MockMetrics{}}}
}

type MockMetrics struct{}

func (m *MockMetrics) RecordPutRequest(status string, duration *time.Time) {
	if duration != nil {
		HT1["puts.current_url.duration"] = time.Since(*duration).Seconds()
	} else {
		switch status {
		case "add":
			HT2["puts.current_url.request.total"] = HT2["puts.current_url.request.total"] + 1
		case "error":
			HT2["puts.current_url.request.error"] = HT2["puts.current_url.request.error"] + 1
		case "bad_request":
			HT2["puts.current_url.request.bad_request"] = HT2["puts.current_url.request.bad_request"] + 1
		}
	}
}

func (m *MockMetrics) RecordGetRequest(status string, duration *time.Time) {
	if duration != nil {
		HT1["gets.current_url.duration"] = time.Since(*duration).Seconds()
	} else {
		switch status {
		case "add":
			HT2["gets.current_url.request.total"] = HT2["gets.current_url.request.total"] + 1
		case "error":
			HT2["gets.current_url.request.error"] = HT2["gets.current_url.request.error"] + 1
		case "bad_request":
			HT2["gets.current_url.request.bad_request"] = HT2["gets.current_url.request.bad_request"] + 1
		}
	}
}
func (m *MockMetrics) RecordPutBackendRequest(status string, duration *time.Time, sizeInBytes float64) {
	if duration != nil {
		HT1["puts.backends.request_duration"] = time.Since(*duration).Seconds()
	} else if sizeInBytes > 0 {
		HT1["puts.backends.request_size_bytes"] = sizeInBytes
	} else {
		switch status {
		case "add":
			HT2["puts.backends.request.total"] = HT2["puts.backends.request.total"] + 1
		case "json":
			HT2["puts.backends.json"] = HT2["puts.backends.json"] + 1
		case "xml":
			HT2["puts.backends.xml"] = HT2["puts.backends.xml"] + 1
		case "invalid_format":
			HT2["puts.backends.invalid_format"] = HT2["puts.backends.invalid_format"] + 1
		case "defines_ttl":
			HT2["puts.backends.defines_ttl"] = HT2["puts.backends.defines_ttl"] + 1
		case "error":
			HT2["puts.backends.request.error"] = HT2["puts.backends.request.error"] + 1
		}
	}
}

func (m *MockMetrics) RecordGetBackendRequest(status string, duration *time.Time) {
	if duration != nil {
		HT1["gets.backends.duration"] = time.Since(*duration).Seconds()
	} else {
		switch status {
		case "add":
			HT2["gets.backends.request.total"] = HT2["gets.backends.request.total"] + 1
		case "error":
			HT2["gets.backends.request.error"] = HT2["gets.backends.request.error"] + 1
		case "bad_request":
			HT2["gets.backends.request.bad_request"] = HT2["gets.backends.request.bad_request"] + 1
		}
	}
}
func (m *MockMetrics) RecordConnectionMetrics(label string) {
	switch label {
	case "add":
		HT1["connections.connections_opened"] = HT1["connections.connections_opened"] + 1
	case "substract":
		HT1["connections.connections_opened"] = HT1["connections.connections_opened"] - 1
	case "accept":
		HT2["connections.connection_error.accept"] = HT2["connections.connection_error.accept"] + 1
	case "close":
		HT2["connections.connection_error.close"] = HT2["connections.connection_error.close"] + 1
	}
}
func (m *MockMetrics) RecordExtraTTLSeconds(aVar float64) {
	HT1["extra_ttl_seconds"] = aVar
}
func (m *MockMetrics) Export(cfg config.Metrics) {
	//
}
