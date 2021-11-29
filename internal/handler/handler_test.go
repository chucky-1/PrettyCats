package handler

import (
	"CatsCrud/internal/request"
	"CatsCrud/internal/service/mock"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"

	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestCatHandler_GetAllCats(t *testing.T) {
	// Setup
	inputJSON := `{}`
	catsJSON := `[{"id":0,"name":""}]`
	e := echo.New()
	e.Validator = &request.CustomValidator{Validator: validator.New()}
	req := httptest.NewRequest(http.MethodGet, "/cats", strings.NewReader(inputJSON))
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	srv := mock.NewService()
	var Handler = NewHandler(srv)

	// Assertions
	if assert.NoError(t, Handler.GetAll(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, catsJSON, strings.Trim(rec.Body.String(), "\n"))
	}
}

func TestCatHandler_CreateCats(t *testing.T) {
	TestTable := []struct {
		name             string
		inputJSON        string
		exceptStatusCode int
	}{
		{
			name: "OK", inputJSON: `{"id":1,"name":"Jon Snow"}`, exceptStatusCode: http.StatusCreated,
		},
		{
			name: "Name is nill", inputJSON: `{"id":1}`, exceptStatusCode: http.StatusBadRequest,
		},
		{
			name: "ID is nill", inputJSON: `{"name":"1"}`, exceptStatusCode: http.StatusBadRequest,
		},
		{
			name: "Params isn't valid", inputJSON: `{"id":"1", "name":"Jon Snow"}`, exceptStatusCode: http.StatusBadRequest,
		},
		{
			name: "name too small", inputJSON: `{"id":1,"name":"Jo"}`, exceptStatusCode: http.StatusBadRequest,
		},
	}

	for _, TestCase := range TestTable {
		TestCase := TestCase // pin!
		t.Run(TestCase.name, func(t *testing.T) {
			e := echo.New()
			e.Validator = &request.CustomValidator{Validator: validator.New()}
			req := httptest.NewRequest(http.MethodPost, "/cats", strings.NewReader(TestCase.inputJSON))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			srv := mock.NewService()
			var Handler = NewHandler(srv)

			if assert.NoError(t, Handler.Create(c)) {
				assert.Equal(t, TestCase.exceptStatusCode, rec.Code)
			}
		})
	}
}

func TestCatHandler_GetCat(t *testing.T) {
	TestTable := []struct {
		name             string
		setParamNames    string
		setParamValues   string
		exceptStatusCode int
	}{
		{
			name: "OK", setParamNames: "id", setParamValues: "1", exceptStatusCode: http.StatusOK,
		},
		{
			name: "ID is nill", setParamNames: "", setParamValues: "", exceptStatusCode: http.StatusBadRequest,
		},
	}

	for _, TestCase := range TestTable {
		TestCase := TestCase // pin!
		t.Run(TestCase.name, func(t *testing.T) {
			e := echo.New()
			e.Validator = &request.CustomValidator{Validator: validator.New()}
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetPath("/cats/:id")
			c.SetParamNames(TestCase.setParamNames)
			c.SetParamValues(TestCase.setParamValues)
			srv := mock.NewService()
			var Handler = NewHandler(srv)

			if assert.NoError(t, Handler.Get(c)) {
				assert.Equal(t, TestCase.exceptStatusCode, rec.Code)
			}
		})
	}
}

func TestCatHandler_UpdateCat(t *testing.T) {
	TestTable := []struct {
		name             string
		setParamNames    string
		setParamValues   string
		inputJSON        string
		exceptStatusCode int
	}{
		{
			name: "OK", setParamNames: "id", setParamValues: "1", inputJSON: `{"name":"Jon Snow"}`, exceptStatusCode: http.StatusOK,
		},
		{
			name: "ID is nill", setParamNames: "", setParamValues: "", inputJSON: `{"name":"Jon Snow"}`, exceptStatusCode: http.StatusBadRequest,
		},
		{
			name: "Name is nill", setParamNames: "id", setParamValues: "1", inputJSON: `{}`, exceptStatusCode: http.StatusBadRequest,
		},
		{
			name: "ID isn't int", setParamNames: "id", setParamValues: "text", inputJSON: `{}`, exceptStatusCode: http.StatusBadRequest,
		},
		{
			name: "Name isn't string", setParamNames: "id", setParamValues: "1", inputJSON: `{"name":1}`, exceptStatusCode: http.StatusBadRequest,
		},
	}

	for _, TestCase := range TestTable {
		TestCase := TestCase // pin!
		t.Run(TestCase.name, func(t *testing.T) {
			e := echo.New()
			e.Validator = &request.CustomValidator{Validator: validator.New()}
			req := httptest.NewRequest(http.MethodPut, "/", strings.NewReader(TestCase.inputJSON))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetPath("/cats/:id")
			c.SetParamNames(TestCase.setParamNames)
			c.SetParamValues(TestCase.setParamValues)
			srv := mock.NewService()
			var Handler = NewHandler(srv)

			if assert.NoError(t, Handler.Update(c)) {
				assert.Equal(t, TestCase.exceptStatusCode, rec.Code)
			}
		})
	}
}

func TestCatHandler_DeleteCat(t *testing.T) {
	TestTable := []struct {
		name             string
		setParamNames    string
		setParamValues   string
		exceptStatusCode int
	}{
		{
			name: "OK", setParamNames: "id", setParamValues: "1", exceptStatusCode: http.StatusOK,
		},
		{
			name: "ID is nill", setParamNames: "", setParamValues: "", exceptStatusCode: http.StatusBadRequest,
		},
		{
			name: "ID isn't int", setParamNames: "id", setParamValues: "text", exceptStatusCode: http.StatusBadRequest,
		},
	}

	for _, TestCase := range TestTable {
		TestCase := TestCase // pin!
		t.Run(TestCase.name, func(t *testing.T) {
			e := echo.New()
			e.Validator = &request.CustomValidator{Validator: validator.New()}
			req := httptest.NewRequest(http.MethodPut, "/", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetPath("/cats/:id")
			c.SetParamNames(TestCase.setParamNames)
			c.SetParamValues(TestCase.setParamValues)
			srv := mock.NewService()
			var Handler = NewHandler(srv)

			if assert.NoError(t, Handler.Delete(c)) {
				assert.Equal(t, TestCase.exceptStatusCode, rec.Code)
			}
		})
	}
}
