package handler

//func TestUserAuthHandler_Restricted(t *testing.T) {
//	e := echo.New()
//	req := httptest.NewRequest(http.MethodPost, "/register", strings.NewReader(`{}`))
//	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
//	rec := httptest.NewRecorder()
//	c := e.NewContext(req, rec)
//
//	var srv service.Auth
//	srv = mock.NewMockCatServ()
//	userHandler = NewUserAuthHandler(srv)
//
//	if assert.NoError(t, userHandler.Restricted(c)) {
//		//assert.Equal(t, 200, rec.Code)
//		//assert.Equal(t, TestCase.exceptBody, strings.Trim(rec.Body.String(), "\n"))
//
//	}
//}