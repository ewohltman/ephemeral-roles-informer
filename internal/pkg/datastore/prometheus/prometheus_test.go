package prometheus

import (
	"context"
	"fmt"
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
	contextTimeout = 10 * time.Second

	queryResponseFile      = "testdata/queryResponse.json"
	queryErrorResponseFile = "testdata/queryErrorResponse.json"
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

	_, err = testQuery(t, queryErrorResponseFile)
	if err == nil {
		t.Errorf("Unexpected success with query error response: %s", err)
	}

	actual, err := testQuery(t, queryResponseFile)
	if err != nil {
		t.Fatalf("Error performing test: %s", err)
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

func testQuery(t *testing.T, responseFile string) ([]int, error) {
	testPrometheusServer := httptest.NewServer(testServerHandler(t, responseFile))
	defer testPrometheusServer.Close()

	prometheusProvider, err := NewProvider(testPrometheusServer.URL)
	if err != nil {
		return nil, fmt.Errorf("error creating new Prometheus provider: %w", err)
	}

	ctx, ctxCancel := context.WithTimeout(context.Background(), contextTimeout)
	defer ctxCancel()

	actual, err := prometheusProvider.ProvideShardServerCounts(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting shard server counts: %w", err)
	}

	return actual, nil
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
