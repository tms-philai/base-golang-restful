package testing

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type TestServer struct {
	router *gin.Engine
	t      *testing.T
}

func NewTestServer(t *testing.T, router *gin.Engine) *TestServer {
	gin.SetMode(gin.TestMode)
	return &TestServer{
		router: router,
		t:      t,
	}
}

func (ts *TestServer) GET(path string, headers map[string]string) *httptest.ResponseRecorder {
	req := httptest.NewRequest("GET", path, nil)
	
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	w := httptest.NewRecorder()
	ts.router.ServeHTTP(w, req)

	return w
}

func (ts *TestServer) POST(path string, body interface{}, headers map[string]string) *httptest.ResponseRecorder {
	jsonBody, err := json.Marshal(body)
	assert.NoError(ts.t, err)

	req := httptest.NewRequest("POST", path, bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	w := httptest.NewRecorder()
	ts.router.ServeHTTP(w, req)

	return w
}

func (ts *TestServer) PUT(path string, body interface{}, headers map[string]string) *httptest.ResponseRecorder {
	jsonBody, err := json.Marshal(body)
	assert.NoError(ts.t, err)

	req := httptest.NewRequest("PUT", path, bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	w := httptest.NewRecorder()
	ts.router.ServeHTTP(w, req)

	return w
}

func (ts *TestServer) DELETE(path string, headers map[string]string) *httptest.ResponseRecorder {
	req := httptest.NewRequest("DELETE", path, nil)
	
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	w := httptest.NewRecorder()
	ts.router.ServeHTTP(w, req)

	return w
}

func (ts *TestServer) AssertStatus(w *httptest.ResponseRecorder, expectedStatus int) {
	assert.Equal(ts.t, expectedStatus, w.Code)
}

func (ts *TestServer) AssertJSON(w *httptest.ResponseRecorder, expected interface{}) {
	var actual interface{}
	err := json.Unmarshal(w.Body.Bytes(), &actual)
	assert.NoError(ts.t, err)
	
	expectedJSON, err := json.Marshal(expected)
	assert.NoError(ts.t, err)
	
	var expectedMap interface{}
	err = json.Unmarshal(expectedJSON, &expectedMap)
	assert.NoError(ts.t, err)
	
	assert.Equal(ts.t, expectedMap, actual)
}

func (ts *TestServer) ParseJSON(w *httptest.ResponseRecorder, dest interface{}) {
	err := json.Unmarshal(w.Body.Bytes(), dest)
	assert.NoError(ts.t, err)
}

type APITestCase struct {
	Name           string
	Method         string
	Path           string
	Body           interface{}
	Headers        map[string]string
	ExpectedStatus int
	ExpectedBody   interface{}
	Setup          func()
	Teardown       func()
}

func (ts *TestServer) RunTestCases(testCases []APITestCase) {
	for _, tc := range testCases {
		ts.t.Run(tc.Name, func(t *testing.T) {
			if tc.Setup != nil {
				tc.Setup()
			}

			var w *httptest.ResponseRecorder

			switch tc.Method {
			case "GET":
				w = ts.GET(tc.Path, tc.Headers)
			case "POST":
				w = ts.POST(tc.Path, tc.Body, tc.Headers)
			case "PUT":
				w = ts.PUT(tc.Path, tc.Body, tc.Headers)
			case "DELETE":
				w = ts.DELETE(tc.Path, tc.Headers)
			default:
				t.Fatalf("Unsupported HTTP method: %s", tc.Method)
			}

			assert.Equal(t, tc.ExpectedStatus, w.Code)

			if tc.ExpectedBody != nil {
				ts.AssertJSON(w, tc.ExpectedBody)
			}

			if tc.Teardown != nil {
				tc.Teardown()
			}
		})
	}
}
