package util

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"runtime"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/render"
)

//Next useful functions are only for testing code

// GetFunctionName return function name, it is useful for mocking.
func GetFunctionName(f interface{}) string {
	fullName := runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
	splitFullName := strings.Split(fullName, ".")
	name := splitFullName[len(splitFullName)-1]
	return strings.Split(name, "-")[0]
}

// RenderToJSON generates a response as a controller response.
func RenderToJSON(u any) string {
	r := render.JSON{Data: u}
	w := httptest.NewRecorder()
	if err := r.Render(w); err != nil {
		panic(err)
	}
	return w.Body.String()
}

// GetTestContext generates a gin.context for test. Useful for testing controller layer
func GetTestContext(params map[string]string, queryString string, body any) (*gin.Context, *httptest.ResponseRecorder, error) {
	gin.SetMode(gin.TestMode)
	r := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(r)
	ctx.Request = &http.Request{
		URL: &url.URL{
			RawQuery: queryString,
		},
	}

	for key, value := range params {
		ctx.AddParam(key, value)
	}

	if body == nil {
		return ctx, r, nil
	}

	if b, err := json.Marshal(body); err != nil {
		return nil, nil, err
	} else {
		ctx.Request.Body = io.NopCloser(bytes.NewBuffer(b))
	}

	return ctx, r, nil
}
