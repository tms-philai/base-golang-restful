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

type TestHelper struct {
	router *gin.Engine
	t      *testing.T
}

func NewTestHelper(t *testing.T, router *gin.Engine) *TestHelper {
	gin.SetMode(gin.TestMode)
	return &TestHelper{
		router: router,
		t:      t,
	}
}

func (th *TestHelper) MakeRequest(method, path string, body interface{}, headers map[string]string) *httptest.ResponseRecorder {
	var req *http.Request

	if body != nil {
		jsonBody, err := json.Marshal(body)
		assert.NoError(th.t, err)
		req = httptest.NewRequest(method, path, bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
	} else {
		req = httptest.NewRequest(method, path, nil)
	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	w := httptest.NewRecorder()
	th.router.ServeHTTP(w, req)

	return w
}

func (th *TestHelper) GET(path string, headers ...map[string]string) *httptest.ResponseRecorder {
	h := make(map[string]string)
	if len(headers) > 0 {
		h = headers[0]
	}
	return th.MakeRequest("GET", path, nil, h)
}

func (th *TestHelper) POST(path string, body interface{}, headers ...map[string]string) *httptest.ResponseRecorder {
	h := make(map[string]string)
	if len(headers) > 0 {
		h = headers[0]
	}
	return th.MakeRequest("POST", path, body, h)
}

func (th *TestHelper) PUT(path string, body interface{}, headers ...map[string]string) *httptest.ResponseRecorder {
	h := make(map[string]string)
	if len(headers) > 0 {
		h = headers[0]
	}
	return th.MakeRequest("PUT", path, body, h)
}

func (th *TestHelper) PATCH(path string, body interface{}, headers ...map[string]string) *httptest.ResponseRecorder {
	h := make(map[string]string)
	if len(headers) > 0 {
		h = headers[0]
	}
	return th.MakeRequest("PATCH", path, body, h)
}

func (th *TestHelper) DELETE(path string, headers ...map[string]string) *httptest.ResponseRecorder {
	h := make(map[string]string)
	if len(headers) > 0 {
		h = headers[0]
	}
	return th.MakeRequest("DELETE", path, nil, h)
}

func (th *TestHelper) AssertStatus(w *httptest.ResponseRecorder, expectedStatus int) {
	assert.Equal(th.t, expectedStatus, w.Code, "Status code mismatch. Body: %s", w.Body.String())
}

func (th *TestHelper) AssertJSON(w *httptest.ResponseRecorder, expected interface{}) {
	var actual interface{}
	err := json.Unmarshal(w.Body.Bytes(), &actual)
	assert.NoError(th.t, err)

	expectedJSON, err := json.Marshal(expected)
	assert.NoError(th.t, err)

	var expectedMap interface{}
	err = json.Unmarshal(expectedJSON, &expectedMap)
	assert.NoError(th.t, err)

	assert.Equal(th.t, expectedMap, actual)
}

func (th *TestHelper) ParseJSON(w *httptest.ResponseRecorder, dest interface{}) {
	err := json.Unmarshal(w.Body.Bytes(), dest)
	assert.NoError(th.t, err, "Failed to parse JSON response: %s", w.Body.String())
}

func (th *TestHelper) AssertContains(w *httptest.ResponseRecorder, substr string) {
	assert.Contains(th.t, w.Body.String(), substr)
}

func (th *TestHelper) AssertNotContains(w *httptest.ResponseRecorder, substr string) {
	assert.NotContains(th.t, w.Body.String(), substr)
}

func (th *TestHelper) GetBody(w *httptest.ResponseRecorder) string {
	return w.Body.String()
}

type TestCase struct {
	Name           string
	Method         string
	Path           string
	Body           interface{}
	Headers        map[string]string
	ExpectedStatus int
	ExpectedBody   interface{}
	CheckResponse  func(*httptest.ResponseRecorder)
	Setup          func()
	Teardown       func()
}

func (th *TestHelper) RunTestCases(testCases []TestCase) {
	for _, tc := range testCases {
		th.t.Run(tc.Name, func(t *testing.T) {
			if tc.Setup != nil {
				tc.Setup()
			}

			w := th.MakeRequest(tc.Method, tc.Path, tc.Body, tc.Headers)

			if tc.ExpectedStatus != 0 {
				assert.Equal(t, tc.ExpectedStatus, w.Code, "Status code mismatch for test: %s. Body: %s", tc.Name, w.Body.String())
			}

			if tc.ExpectedBody != nil {
				th.AssertJSON(w, tc.ExpectedBody)
			}

			if tc.CheckResponse != nil {
				tc.CheckResponse(w)
			}

			if tc.Teardown != nil {
				tc.Teardown()
			}
		})
	}
}

type AuthHelper struct {
	th    *TestHelper
	token string
}

func (th *TestHelper) Auth() *AuthHelper {
	return &AuthHelper{th: th}
}

func (ah *AuthHelper) Login(email, password string) string {
	w := ah.th.POST("/api/v1/auth/login", map[string]string{
		"email":    email,
		"password": password,
	})

	ah.th.AssertStatus(w, http.StatusOK)

	var response map[string]interface{}
	ah.th.ParseJSON(w, &response)

	data, ok := response["data"].(map[string]interface{})
	assert.True(ah.th.t, ok, "Failed to extract data from login response")

	token, ok := data["access_token"].(string)
	assert.True(ah.th.t, ok, "Failed to extract access_token from login response")

	ah.token = token
	return token
}

func (ah *AuthHelper) GetAuthHeader() map[string]string {
	return map[string]string{
		"Authorization": "Bearer " + ah.token,
	}
}

func (ah *AuthHelper) GET(path string) *httptest.ResponseRecorder {
	return ah.th.GET(path, ah.GetAuthHeader())
}

func (ah *AuthHelper) POST(path string, body interface{}) *httptest.ResponseRecorder {
	return ah.th.POST(path, body, ah.GetAuthHeader())
}

func (ah *AuthHelper) PUT(path string, body interface{}) *httptest.ResponseRecorder {
	return ah.th.PUT(path, body, ah.GetAuthHeader())
}

func (ah *AuthHelper) DELETE(path string) *httptest.ResponseRecorder {
	return ah.th.DELETE(path, ah.GetAuthHeader())
}

func CreateTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	return r
}

func AssertJSONField(t *testing.T, w *httptest.ResponseRecorder, field string, expectedValue interface{}) {
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	actualValue, exists := response[field]
	assert.True(t, exists, "Field %s not found in response", field)
	assert.Equal(t, expectedValue, actualValue, "Value mismatch for field %s", field)
}

func AssertErrorResponse(t *testing.T, w *httptest.ResponseRecorder, expectedStatus int, expectedErrorMsg string) {
	assert.Equal(t, expectedStatus, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	if expectedErrorMsg != "" {
		message, exists := response["message"]
		assert.True(t, exists, "Error message not found in response")
		assert.Contains(t, message, expectedErrorMsg)
	}
}
