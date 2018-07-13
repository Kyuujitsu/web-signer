package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

type DummySigner struct{}

func (s *DummySigner) Sign(datahex, privhex string) (string, error) {
	return fmt.Sprintf("%s|%s", datahex, privhex), nil
}

func TestGetSignature(t *testing.T) {
	server := NewSignServer(&DummySigner{})

	t.Run("it returns signature based on data provided", func(t *testing.T) {
		datahex := "123"
		keyhex := "321"
		request := newGetSignRequest(datahex, keyhex)
		response := httptest.NewRecorder()

		expectedResult := "123|321"

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusOK)
		assertResponseBody(t, response.Body.String(), expectedResult)
	})

	t.Run("it returns Bad Request error when no or incorrect data provided", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/sign", nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusBadRequest)
	})
}

func newGetSignRequest(data, private string) *http.Request {
	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/sign?data=%s&private=%s", data, private), nil)
	return req
}

func assertStatus(t *testing.T, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf("didn't get expected status, got %d, expected %d", got, want)
	}
}

func assertResponseBody(t *testing.T, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf("did not provide a correct response, got %s, want %s", got, want)
	}
}
