package util

import (
	"net/http/httptest"
	"reflect"
	"runtime"
	"strings"

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
