package handler_test

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"workshop/internal/handler"

	"workshop/internal/api"
	"workshop/internal/api/mocks"
)

func TestHandler_Hello(t *testing.T) {
	tests := []struct {
		name     string
		joke     *api.JokeResponse
		err      error
		codeWant int
		bodyWant string
	}{
		{
			name:     "simple test",
			joke:     &api.JokeResponse{Joke: "some joke"},
			err:      nil,
			codeWant: http.StatusOK,
			bodyWant: "some joke",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			apiMock := &mocks.Client{}
			apiMock.On("GetJoke").Return(tt.joke, tt.err)
			h := handler.NewHandler(apiMock)

			req, _ := http.NewRequest("GET", "/hello", nil)
			rr := httptest.NewRecorder()
			handlerFunc := http.HandlerFunc(h.Hello)
			handlerFunc.ServeHTTP(rr, req)

			gotRaw, _ := ioutil.ReadAll(rr.Body)
			got := string(gotRaw)
			if got != tt.bodyWant {
				t.Errorf("wrong responde body %s want %s", got, tt.bodyWant)
			}

			if status := rr.Result().StatusCode; status != tt.codeWant {
				t.Errorf("wrong status code %d want %d", status, tt.codeWant)
			}
		})
	}
}
