package server

import (
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/tsedgwick/hash-api/api"
)

func TestV1HashHandler(t *testing.T) {
	tests := []struct {
		name               string
		method             string
		expected           string
		input              func() io.Reader
		expectedStatusCode int
	}{
		{
			name:     "Success",
			method:   http.MethodPost,
			expected: "token",
			input: func() io.Reader {
				form := url.Values{}
				form.Add("password", "angryMonkey")
				return strings.NewReader(form.Encode())
			},
			expectedStatusCode: http.StatusOK,
		},
		{
			name:   "Wrong method",
			method: http.MethodGet,
			input: func() io.Reader {
				return nil
			},
			expectedStatusCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			s := New(":8080", &mockClient{
				token: "token",
				key:   "123",
			})
			req, err := http.NewRequest(tt.method, "/v1/hash", tt.input())
			req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
			if err != nil {
				t.Fatal(err)
			}
			rr := httptest.NewRecorder()
			http.HandlerFunc(s.v1HashHandler).ServeHTTP(rr, req)
			actualCode := rr.Code
			if actualCode != tt.expectedStatusCode {
				t.Errorf("failed at %s : wrong error code : expected %d : actual : %d", tt.name, tt.expectedStatusCode, actualCode)
			}

			actual, err := ioutil.ReadAll(rr.Body)
			if err != nil {
				t.Fatalf("Error occurred :%v", err)
			}
			if string(actual) != tt.expected {
				t.Errorf("failed at %s : expected %v : actual : %v", tt.name, tt.expected, actual)
			}
		})
	}
}

func TestV2HashHandler(t *testing.T) {
	tests := []struct {
		name               string
		method             string
		expected           string
		input              func() io.Reader
		expectedStatusCode int
	}{
		{
			name:   "Success",
			method: http.MethodPost,
			input: func() io.Reader {
				form := url.Values{}
				form.Add("password", "angryMonkey")
				return strings.NewReader(form.Encode())
			},
			expected:           "123",
			expectedStatusCode: http.StatusOK,
		},
		{
			name:   "Wrong method",
			method: http.MethodGet,
			input: func() io.Reader {
				return nil
			},
			expectedStatusCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			s := New(":8080", &mockClient{
				token: "token",
				key:   "123",
			})
			req, err := http.NewRequest(tt.method, "/v1/hash", tt.input())
			req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
			if err != nil {
				t.Fatal(err)
			}
			rr := httptest.NewRecorder()
			http.HandlerFunc(s.v2HashHandler).ServeHTTP(rr, req)
			actualCode := rr.Code
			if actualCode != tt.expectedStatusCode {
				t.Errorf("failed at %s : wrong error code : expected %d : actual : %d", tt.name, tt.expectedStatusCode, actualCode)
			}

			actual, err := ioutil.ReadAll(rr.Body)
			if err != nil {
				t.Fatalf("Error occurred :%v", err)
			}
			if string(actual) != tt.expected {
				t.Errorf("failed at %s : expected %v : actual : %v", tt.name, tt.expected, actual)
			}
		})
	}
}

func TestV3HashHandler(t *testing.T) {
	tests := []struct {
		name               string
		method             string
		key                string
		expected           string
		client             api.Client
		expectedStatusCode int
	}{
		{
			name:     "Success",
			method:   http.MethodGet,
			key:      "123",
			expected: "token",
			client: &mockClient{
				token: "token",
				key:   "123",
			},
			expectedStatusCode: http.StatusOK,
		},
		{
			name:               "Wrong method",
			method:             http.MethodPost,
			expectedStatusCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			s := New(":8080", tt.client)
			req, err := http.NewRequest(tt.method, "/v3/hash/"+tt.key, nil)
			if err != nil {
				t.Fatal(err)
			}
			rr := httptest.NewRecorder()
			http.HandlerFunc(s.v3HashHandler).ServeHTTP(rr, req)
			actualCode := rr.Code
			if actualCode != tt.expectedStatusCode {
				t.Errorf("failed at %s : wrong error code : expected %d : actual : %d", tt.name, tt.expectedStatusCode, actualCode)
			}

			actual, err := ioutil.ReadAll(rr.Body)
			if err != nil {
				t.Fatalf("Error occurred :%v", err)
			}
			if string(actual) != tt.expected {
				t.Errorf("failed at %s : expected %v : actual : %v", tt.name, tt.expected, actual)
			}

		})
	}
}

func TestShutdownHandler(t *testing.T) {
	tests := []struct {
		name               string
		method             string
		expected           string
		expectedStatusCode int
	}{
		{
			name:               "Success",
			method:             http.MethodPost,
			expectedStatusCode: http.StatusOK,
		},
		{
			name:               "Wrong method",
			method:             http.MethodGet,
			expectedStatusCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			s := New(":8080", &mockClient{
				token: "token",
				key:   "123",
			})
			req, err := http.NewRequest(tt.method, "/shutdown", nil)

			if err != nil {
				t.Fatal(err)
			}
			rr := httptest.NewRecorder()
			http.HandlerFunc(s.shutdownHandler).ServeHTTP(rr, req)
			actualCode := rr.Code
			if actualCode != tt.expectedStatusCode {
				t.Errorf("failed at %s : wrong error code : expected %d : actual : %d", tt.name, tt.expectedStatusCode, actualCode)
			}
		})
	}
}

type mockClient struct {
	key   string
	token string
}

func (m *mockClient) Encode(input []byte) string {
	return m.token
}
func (m *mockClient) Tokenize(input []byte) string {
	return m.token

}
func (m *mockClient) Token(key string) string {
	return m.token

}
func (m *mockClient) Save(input []byte) string {
	return m.key

}
