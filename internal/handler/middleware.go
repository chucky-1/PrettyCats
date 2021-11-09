package handler

import (
	"CatsCrud/internal/service"
	"fmt"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

func (h *UserAuthHandler) Restricted (c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	fmt.Println(user.Raw)
	fmt.Println(user.Method)
	fmt.Println(user.Header)
	fmt.Println(user.Claims)
	fmt.Println(user.Signature)
	fmt.Println(user.Valid)
	claims := user.Claims.(*service.JwtCustomClaims)
	name := claims.Name
	id := claims.ID
	idStr := strconv.Itoa(id)
	return c.String(http.StatusOK, "Welcome " + name + idStr)
}
