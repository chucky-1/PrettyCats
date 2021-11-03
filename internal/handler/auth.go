package handler

import (
	"CatsCrud/internal/models"
	"CatsCrud/internal/service"
	"encoding/json"
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
)

type UserAuthHandler struct {
	src service.Auth
}

func NewUserAuthHandler(src service.Auth) *UserAuthHandler {
	return &UserAuthHandler{src: src}
}

func (h *UserAuthHandler) SignUp(c echo.Context) error {
	fmt.Println("SignUp start...")
	var user models.User

	err := json.NewDecoder(c.Request().Body).Decode(&user)
	if err != nil {
		return c.JSON(http.StatusBadRequest, new(models.User))
	}

	fmt.Println("Decode successful")
	fmt.Println(user)

	// Проверка, что все поля заполнены
	if user.Username == "" || user.Name == "" || user.Password == "" {
		return c.JSON(http.StatusBadRequest, new(models.User))
	}

	fmt.Println("Params corrected")

	id, err := h.src.CreateUserServ(user)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	fmt.Println("Вернулся id: ", id)
	
	return nil
}
