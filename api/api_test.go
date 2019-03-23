package api

import (
	"github.com/MultiplePingChecker/api"
	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
	//"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

var (
	pingForm = "server=google.com&count=100"
	pingJSON = `{"hostname":"google.com","count":100}
`
)
type 
func TestCreatePing(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/ping", strings.NewReader(pingForm))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if assert.NoError(t, api.CreatePing(c)) {
		assert.Equal(t, http.StatusCreated, rec.Code)
		assert.Equal(t, pingJSON, rec.Body.String())
	}
}
func TestGetPingStatus(t *testing.T) {
	e := echo.New()
	req := httptest.NewReqeust(http.MethodGet, "/ping", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/:hostname")
	c.SetParamName("hostname")
	c.SetParamValues("google.com")
}
