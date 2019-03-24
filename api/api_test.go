package api

import (
	"github.com/MultiplePingChecker/api"
	"github.com/MultiplePingChecker/temp_db"
	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
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

type dbMock struct {
	mock.Mock
}

func (d *dbMock) InsertPingHost(key string, value *temp_db.PingHostEntry) int {
	args := d.Called(key, value)
	return args.Int(0)
}
func (d *dbMock) DeletePingHost(key string) int {
	args := d.Called(key)
	return args.Int(0)

}
func (d *dbMock) SearchHost(key string) (*temp_db.PingHostEntry, bool) {
	args := d.Called(key)
	return args.Get(0).(*temp_db.PingHostEntry), args.Bool(1)
}
func (d *dbMock) GetPingHost() *map[string]*temp_db.PingHostEntry {
	args := d.Called()
	return args.Get(0).(*map[string]*temp_db.PingHostEntry)
}
func TestCreatePing(t *testing.T) {
	e := echo.New()
	DBMock := dbMock{}
	DBMock.On("InsertPingHost", "google.com", mock.Anything).Return(1)
	a := api.Api{&DBMock}
	req := httptest.NewRequest(http.MethodPost, "/ping", strings.NewReader(pingForm))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if assert.NoError(t, a.CreatePing(c)) {
		assert.Equal(t, http.StatusCreated, rec.Code)
		assert.Equal(t, pingJSON, rec.Body.String())
	}
	DBMock.AssertExpectations(t)
	mock.AssertExpectationsForObjects(t, &DBMock)
}
func TestGetPing(t *testing.T) {
	e := echo.New()
	DBMock := dbMock{}
	DBMock.On("GetPingHost").Return(&map[string]*temp_db.PingHostEntry{})
	a := api.Api{&DBMock}
	req := httptest.NewRequest(http.MethodGet, "/ping", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if assert.NoError(t, a.GetPing(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "[]\n", rec.Body.String())
	}
}
func TestGetPingStatus(t *testing.T) {
	e := echo.New()
	DBMock := dbMock{}
	DBMock.On("SearchHost", "google.com").Return(&temp_db.PingHostEntry{}, false)
	a := api.Api{&DBMock}
	req := httptest.NewRequest(http.MethodGet, "/ping", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/:hostname")
	c.SetParamNames("hostname")
	c.SetParamValues("google.com")

	if assert.NoError(t, a.DeletePing(c)) {
		assert.Equal(t, http.StatusNotFound, rec.Code)
	}
}
func TestDeletePing(t *testing.T) {
	e := echo.New()
	DBMock := dbMock{}
	DBMock.On("SearchHost", "google.com").Return(&temp_db.PingHostEntry{}, false)
	a := api.Api{&DBMock}
	req := httptest.NewRequest(http.MethodDelete, "/ping", strings.NewReader(pingForm))
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/:hostname")
	c.SetParamNames("hostname")
	c.SetParamValues("google.com")

	if assert.NoError(t, a.DeletePing(c)) {
		assert.Equal(t, http.StatusNotFound, rec.Code)
	}
}
