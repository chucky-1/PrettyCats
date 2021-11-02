package handler

import (
	"CatsCrud/internal/models"
	"CatsCrud/internal/service"
	"encoding/json"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

type CatHandler struct {
	src service.Service
}

func NewCatHandler(srv service.Service) *CatHandler {
	return &CatHandler{src: srv}
}

func(h *CatHandler) GetAllCats(c echo.Context) error {
	allcats, err :=  h.src.GetAllCatsServ()
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, allcats)
}

func (h *CatHandler) CreateCats(c echo.Context) error {
	cats := new(models.Cats)
	err := json.NewDecoder(c.Request().Body).Decode(cats)
	if err != nil {
		return c.JSON(http.StatusBadRequest, new(models.Cats))
	}
	if cats.Name == "" {
		return c.JSON(http.StatusBadRequest, new(models.Cats))
	}
	if cats.ID == 0 {
		return c.JSON(http.StatusBadRequest, new(models.Cats))
	}

	cat, err := h.src.CreateCatsServ(*cats)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusCreated, cat)
}

func (h *CatHandler) GetCat(c echo.Context) error {
	// Достаём ID
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, new(models.Cats))
	}

	cat, err := h.src.GetCatServ(id)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, cat)
}

func (h *CatHandler) UpdateCat(c echo.Context) error {
	cats := new(models.Cats)

	// Достаём ID
	id := c.Param("id")
	// Проверка что id передан
	if id == "" {
		return c.JSON(http.StatusBadRequest, new(models.Cats))
	}
	// Проверка что id можно перевезти в int
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, new(models.Cats))
	}
	// Присваиваем id модели
	cats.ID = int32(idInt)

	err = json.NewDecoder(c.Request().Body).Decode(&cats)
	if err != nil {
		return c.JSON(http.StatusBadRequest, new(models.Cats))
	}
	// Проверка что name передан
	if cats.Name == "" {
		return c.JSON(http.StatusBadRequest, new(models.Cats))
	}

	cat, err := h.src.UpdateCatServ(id, *cats)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, cat)
}

func (h *CatHandler) DeleteCat(c echo.Context) error {
	// Достаём ID
	id := c.Param("id")
	// Проверка что id передан
	if id == "" {
		return c.JSON(http.StatusBadRequest, new(models.Cats))
	}
	// Проверка что id можно перевезти в int
	_, err := strconv.Atoi(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, new(models.Cats))
	}

	cat, err := h.src.DeleteCatServ(id)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, cat)
}
