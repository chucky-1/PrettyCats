package handler

import (
	"CatsCrud/internal/service"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)


// @Summary Restricted
// @Security ApiKeyAuth
// @Description example closed page
// @Produce json
// @Success 200 {string} string
// @Failure 400 {object} string
// @Router /restrict [get]
func (h *UserAuthHandler) Restricted (c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*service.JwtCustomClaims)
	name := claims.Name
	id := claims.ID
	idStr := strconv.Itoa(id)
	return c.String(http.StatusOK, "Welcome " + name + idStr)
}
