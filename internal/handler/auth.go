package handler

import (
	"CatsCrud/internal/models"
	"CatsCrud/internal/service"
	"encoding/json"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

type UserAuthHandler struct {
	src service.Auth
}

func NewUserAuthHandler(src service.Auth) *UserAuthHandler {
	return &UserAuthHandler{src: src}
}

func (h *UserAuthHandler) SignUp(c echo.Context) error {
	var user models.User

	err := json.NewDecoder(c.Request().Body).Decode(&user)
	if err != nil {
		return c.JSON(http.StatusBadRequest, new(models.User))
	}

	// Проверка, что все поля заполнены
	if user.Username == "" || user.Name == "" || user.Password == "" {
		return c.JSON(http.StatusBadRequest, new(models.User))
	}

	id, err := h.src.CreateUserServ(user)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, id)
}

type SignInInput struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (h *UserAuthHandler) SignIn(c echo.Context) error {
	var input SignInInput

	err := json.NewDecoder(c.Request().Body).Decode(&input)
	if err != nil {
		return c.JSON(http.StatusBadRequest, new(models.User))
	}

	// Проверка, что все поля заполнены
	if input.Username == "" || input.Password == "" {
		return c.JSON(http.StatusBadRequest, new(models.User))
	}

	token, err := h.src.GenerateToken(input.Username, input.Password)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, token)
}

func (h *UserAuthHandler) Restricted(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*service.JwtCustomClaims)
	name := claims.Name
	id := claims.ID
	idStr := strconv.Itoa(id)
	return c.String(http.StatusOK, "Welcome " + name + idStr)
}
