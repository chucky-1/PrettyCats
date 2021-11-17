package handler

import (
	"CatsCrud/internal/models"
	"CatsCrud/internal/service"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type UserAuthHandler struct {
	src service.Auth
}

func NewUserAuthHandler(src service.Auth) *UserAuthHandler {
	return &UserAuthHandler{src: src}
}

// @Summary SignUp
// @Tags auth
// @Description decode params and send it in service for create account
// @Accept json
// @Produce json
// @Param user body models.User true "user"
// @Success 200 {integer} integer
// @Failure 400 {object} models.User
// @Failure 500 {object} models.User
// @Router /register [post]
func (h *UserAuthHandler) SignUp(c echo.Context) error {
	user := new(models.User)

	err := c.Bind(user)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	if err = c.Validate(user); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	id, err := h.src.CreateUserServ(*user)
	if err != nil {
		log.Error(err)
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, id)
}

type SignInInput struct {
	Username string `json:"username" validate:"required,min=3"`
	Password string `json:"password" validate:"required,min=6"`
}

// @Summary SignIn
// @Tags auth
// @Description decode params and send them in service for generate token
// @Accept json
// @Produce json
// @Param input body SignInInput true "input"
// @Success 200 {string} string "token"
// @Failure 400 {object} models.User
// @Failure 500 {object} models.User
// @Router /login [post]
func (h *UserAuthHandler) SignIn(c echo.Context) error {
	input := new(SignInInput)

	err := c.Bind(input)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	if err = c.Validate(input); err != nil {
		return c.JSON(http.StatusBadRequest,err)
	}

	token, err := h.src.GenerateToken(input.Username, input.Password)
	if err != nil {
		log.Error(err)
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, token)
}
