package bind

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type payload struct {
	Name string `json:"name"`
}

func newContext(body string) echo.Context {
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	return e.NewContext(req, httptest.NewRecorder())
}

func TestBind_Valid(t *testing.T) {
	c := newContext(`{"name":"ivan"}`)

	var p payload
	require.NoError(t, Bind(c, &p))
	assert.Equal(t, "ivan", p.Name)
}

func TestBind_UnknownFieldRejected(t *testing.T) {
	c := newContext(`{"name":"ivan","hacker":true}`)

	var p payload
	err := Bind(c, &p)
	assert.Error(t, err, "unknown fields must be rejected (DisallowUnknownFields)")
}

func TestBind_MalformedJSON(t *testing.T) {
	c := newContext(`{not json`)

	var p payload
	assert.Error(t, Bind(c, &p))
}
