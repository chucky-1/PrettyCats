package handler

import (
	"CatsCrud/internal/request"
	"CatsCrud/internal/service"
	"CatsCrud/internal/service/mock"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"

	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestAuthHandler_SignUp(t *testing.T) {
	TestTable := []struct {
		name             string
		inputJSON        string
		exceptStatusCode int
	}{
		{
			name: "OK", inputJSON: `{"name":"Jon Snow","username":"Jonny","password":"Jon1215d"}`, exceptStatusCode: 200,
		},
		{
			name: "name is missing", inputJSON: `{"username":"Jonny","password":"Jon12@15d"}`, exceptStatusCode: http.StatusBadRequest,
		},
		{
			name: "username is missing", inputJSON: `{"name":"Jon Snow","password":"Jon12@15d"}`, exceptStatusCode: http.StatusBadRequest,
		},
		{
			name: "password is missing", inputJSON: `{"name":"Jon Snow","username":"Jonny"}`, exceptStatusCode: http.StatusBadRequest,
		},
	}

	for _, TestCase := range TestTable {
		TestCase := TestCase // pin!
		t.Run(TestCase.name, func(t *testing.T) {
			e := echo.New()
			e.Validator = &request.CustomValidator{Validator: validator.New()}
			req := httptest.NewRequest(http.MethodPost, "/register", strings.NewReader(TestCase.inputJSON))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			var srv service.Auth = mock.NewService()
			var userHandler = NewAuthHandler(srv)

			if assert.NoError(t, userHandler.SignUp(c)) {
				assert.Equal(t, TestCase.exceptStatusCode, rec.Code)
			}
		})
	}
}

func TestAuthHandler_SignIn(t *testing.T) {
	TestTable := []struct {
		name             string
		inputJSON        string
		exceptBody       string
		exceptStatusCode int
	}{
		{
			name: "OK", inputJSON: `{"username":"Jonny","password":"Jon12@15d"}`, exceptStatusCode: http.StatusOK,
		},
		{
			name: "Username isn't in", inputJSON: `{"password":"Jon12@15d"}`, exceptStatusCode: http.StatusBadRequest,
		},
		{
			name: "Password isn't in", inputJSON: `{"username":"Jonny"}`, exceptStatusCode: http.StatusBadRequest,
		},
	}

	for _, TestCase := range TestTable {
		TestCase := TestCase // pin!
		t.Run(TestCase.name, func(t *testing.T) {
			e := echo.New()
			e.Validator = &request.CustomValidator{Validator: validator.New()}
			req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(TestCase.inputJSON))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			var srv service.Auth = mock.NewService()
			var userHandler = NewAuthHandler(srv)

			if assert.NoError(t, userHandler.SignIn(c)) {
				assert.Equal(t, TestCase.exceptStatusCode, rec.Code)
			}
		})
	}
}
