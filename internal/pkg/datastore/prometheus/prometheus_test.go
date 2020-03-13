package prometheus

import (
	"context"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"reflect"
	"testing"
	"time"
)

const (
	contextTimeout    = 10 * time.Second
	queryResponseFile = "testdata/queryResponse.json"
)

func TestNew(t *testing.T) {
	testPrometheusServer := httptest.NewServer(testServerHandler(t, queryResponseFile))
	defer testPrometheusServer.Close()

	_, err := NewProvider(testPrometheusServer.URL)
	if err != nil {
		t.Errorf("Error creating new Prometheus provider: %s", err)
	}
}

func TestProvider_ProvideShardServerCounts(t *testing.T) {
	_, err := NewProvider("\\http://")
	if err == nil {
		t.Errorf("Unexpected success with invalid URL")
	}

	testPrometheusServer := httptest.NewServer(testServerHandler(t, queryResponseFile))
	defer testPrometheusServer.Close()

	prometheusProvider, err := NewProvider(testPrometheusServer.URL)
	if err != nil {
		t.Fatalf("Error creating new Prometheus provider: %s", err)
	}

	ctx, ctxCancel := context.WithTimeout(context.Background(), contextTimeout)
	defer ctxCancel()

	actual, err := prometheusProvider.ProvideShardServerCounts(ctx)
	if err != nil {
		t.Fatalf("Error getting shard server counts: %s", err)
	}

	expected := []int{31, 28, 27, 23, 26, 19, 24, 23, 17, 28}

	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf(
			"Unexpected result. Got: %v, Expected: %v",
			actual,
			expected,
		)
	}
}

func testServerHandler(t *testing.T, responseFile string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			closeErr := r.Body.Close()
			if closeErr != nil {
				t.Errorf("Error closing request body: %s", closeErr)
			}
		}()

		_, err := io.Copy(ioutil.Discard, r.Body)
		if err != nil {
			t.Errorf("Error reading request body: %s", err)

			writeErrorResponse(t, w, http.StatusInternalServerError, err)

			return
		}

		respBody, err := ioutil.ReadFile(filepath.Clean(responseFile))
		if err != nil {
			t.Errorf("Error reading file %s: %s", queryResponseFile, err)

			writeErrorResponse(t, w, http.StatusInternalServerError, err)

			return
		}

		_, err = w.Write(respBody)
		if err != nil {
			t.Errorf("Error writing response body: %s", err)
		}
	}
}

func writeErrorResponse(t *testing.T, w http.ResponseWriter, status int, err error) {
	w.WriteHeader(status)

	_, writeErr := w.Write([]byte(err.Error()))
	if writeErr != nil {
		t.Errorf("Error writing error response body: %s", err)
	}
}
