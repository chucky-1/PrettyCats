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

var catHandler *CatHandler

func TestCatHandler_GetAllCats(t *testing.T) {
	// Setup
	inputJSON := `{}`
	catsJSON := `[{"id":0,"name":""}]`
	e := echo.New()
	e.Validator = &request.CustomValidator{Validator: validator.New()}
	req := httptest.NewRequest(http.MethodGet, "/cats", strings.NewReader(inputJSON))
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	srv := mock.NewMockCatServ()
	catHandler = NewCatHandler(srv)

	// Assertions
	if assert.NoError(t, catHandler.GetAllCats(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, catsJSON, strings.Trim(rec.Body.String(), "\n"))
	}
}

func TestCatHandler_CreateCats(t *testing.T) {
	TestTable := []struct{
		name string
		inputJson string
		exceptStatusCode int
		exceptBody string
	}{
		{
			name: "OK",
			inputJson: `{"id":1,"name":"Jon Snow"}`,
			exceptStatusCode: http.StatusCreated,
			exceptBody: `{"id":1,"name":"Jon Snow"}`,
		},
		{
			name: "Name is nill",
			inputJson: `{"id":1}`,
			exceptStatusCode: http.StatusBadRequest,
			exceptBody: `{"id":0,"name":""}`,
		},
		{
			name: "ID is nill",
			inputJson: `{"name":"1"}`,
			exceptStatusCode: http.StatusBadRequest,
			exceptBody: `{"id":0,"name":""}`,
		},
		{
			name: "Params isn't valid",
			inputJson: `{"id":"1", "name":"Jon Snow"}`,
			exceptStatusCode: http.StatusBadRequest,
			exceptBody: `{"id":0,"name":""}`,
		},
		{
			name: "name too small",
			inputJson: `{"id":1,"name":"Jo"}`,
			exceptStatusCode: http.StatusBadRequest,
			exceptBody: `{"id":0,"name":""}`,
		},
	}

	for _, TestCase := range TestTable {
		t.Run(TestCase.name, func(t *testing.T) {
			e := echo.New()
			e.Validator = &request.CustomValidator{Validator: validator.New()}
			req := httptest.NewRequest(http.MethodPost, "/cats", strings.NewReader(TestCase.inputJson))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			srv := mock.NewMockCatServ()
			catHandler = NewCatHandler(srv)

			if assert.NoError(t, catHandler.CreateCats(c)) {
				assert.Equal(t, TestCase.exceptStatusCode, rec.Code)
				assert.Equal(t, TestCase.exceptBody, strings.Trim(rec.Body.String(), "\n"))
			}
		})
	}
}

func TestCatHandler_GetCat(t *testing.T) {
	TestTable := []struct{
		name string
		setParamNames string
		setParamValues string
		exceptStatusCode int
		exceptBody string
	}{
		{
			name: "OK",
			setParamNames: "id",
			setParamValues: "1",
			exceptStatusCode: http.StatusOK,
			exceptBody: `{"id":1,"name":"Jon Snow"}`,
		},
		{
			name: "ID is nill",
			setParamNames: "",
			setParamValues: "",
			exceptStatusCode: http.StatusBadRequest,
			exceptBody: `{"id":0,"name":""}`,
		},
	}

	for _, TestCase := range TestTable {
		t.Run(TestCase.name, func(t *testing.T) {
			e := echo.New()
			e.Validator = &request.CustomValidator{Validator: validator.New()}
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetPath("/cats/:id")
			c.SetParamNames(TestCase.setParamNames)
			c.SetParamValues(TestCase.setParamValues)
			srv := mock.NewMockCatServ()
			catHandler = NewCatHandler(srv)

			if assert.NoError(t, catHandler.GetCat(c)) {
				assert.Equal(t, TestCase.exceptStatusCode, rec.Code)
				assert.Equal(t, TestCase.exceptBody, strings.Trim(rec.Body.String(), "\n"))
			}
		})
	}
}

func TestCatHandler_UpdateCat(t *testing.T) {
	TestTable := []struct{
		name string
		setParamNames string
		setParamValues string
		inputJson string
		exceptStatusCode int
		exceptBody string
	}{
		{
			name: "OK",
			setParamNames: "id",
			setParamValues: "1",
			inputJson: `{"name":"Jon Snow"}`,
			exceptStatusCode: http.StatusOK,
			exceptBody: `{"id":1,"name":"Jon Snow"}`,
		},
		{
			name: "ID is nill",
			setParamNames: "",
			setParamValues: "",
			inputJson: `{"name":"Jon Snow"}`,
			exceptStatusCode: http.StatusBadRequest,
			exceptBody: `{"id":0,"name":""}`,
		},
		{
			name: "Name is nill",
			setParamNames: "id",
			setParamValues: "1",
			inputJson: `{}`,
			exceptStatusCode: http.StatusBadRequest,
			exceptBody: `{"id":0,"name":""}`,
		},
		{
			name: "ID isn't int",
			setParamNames: "id",
			setParamValues: "text",
			inputJson: `{}`,
			exceptStatusCode: http.StatusBadRequest,
			exceptBody: `{"id":0,"name":""}`,
		},
		{
			name: "Name isn't string",
			setParamNames: "id",
			setParamValues: "1",
			inputJson: `{"name":1}`,
			exceptStatusCode: http.StatusBadRequest,
			exceptBody: `{"id":0,"name":""}`,
		},
	}

	for _, TestCase := range TestTable {
		t.Run(TestCase.name, func(t *testing.T) {
			e := echo.New()
			e.Validator = &request.CustomValidator{Validator: validator.New()}
			req := httptest.NewRequest(http.MethodPut, "/", strings.NewReader(TestCase.inputJson))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetPath("/cats/:id")
			c.SetParamNames(TestCase.setParamNames)
			c.SetParamValues(TestCase.setParamValues)
			srv := mock.NewMockCatServ()
			catHandler = NewCatHandler(srv)

			if assert.NoError(t, catHandler.UpdateCat(c)) {
				assert.Equal(t, TestCase.exceptStatusCode, rec.Code)
				assert.Equal(t, TestCase.exceptBody, strings.Trim(rec.Body.String(), "\n"))
			}
		})
	}
}

func TestCatHandler_DeleteCat(t *testing.T) {
	TestTable := []struct{
		name string
		setParamNames string
		setParamValues string
		exceptStatusCode int
		exceptBody string
	}{
		{
			name: "OK",
			setParamNames: "id",
			setParamValues: "1",
			exceptStatusCode: http.StatusOK,
			exceptBody: `{"id":1,"name":"Jon Snow"}`,
		},
		{
			name: "ID is nill",
			setParamNames: "",
			setParamValues: "",
			exceptStatusCode: http.StatusBadRequest,
			exceptBody: `{"id":0,"name":""}`,
		},
		{
			name: "ID isn't int",
			setParamNames: "id",
			setParamValues: "text",
			exceptStatusCode: http.StatusBadRequest,
			exceptBody: `{"id":0,"name":""}`,
		},
	}

	for _, TestCase := range TestTable {
		t.Run(TestCase.name, func(t *testing.T) {
			e := echo.New()
			e.Validator = &request.CustomValidator{Validator: validator.New()}
			req := httptest.NewRequest(http.MethodPut, "/", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetPath("/cats/:id")
			c.SetParamNames(TestCase.setParamNames)
			c.SetParamValues(TestCase.setParamValues)
			srv := mock.NewMockCatServ()
			catHandler = NewCatHandler(srv)

			if assert.NoError(t, catHandler.DeleteCat(c)) {
				assert.Equal(t, TestCase.exceptStatusCode, rec.Code)
				assert.Equal(t, TestCase.exceptBody, strings.Trim(rec.Body.String(), "\n"))
			}
		})
	}
}