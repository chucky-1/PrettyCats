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

type partOfCats struct {
	ID   int32  `json:"id" bson:"id" validate:"required,numeric,gt=0"`
}

// @Summary GetAllCats
// @Tags Cats
// @Description collect all cats in array
// @Produce json
// @Success 200 {array} models.Cats
// @Router /cats [get]
func(h *CatHandler) GetAllCats(c echo.Context) error {
	allcats, err :=  h.src.GetAllCatsServ()
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, allcats)
}

// @Summary CreateCats
// @Tags Cats
// @Description create cat
// @Accept json
// @Produce json
// @Param cats body models.Cats true "cats"
// @Success 201 {object} models.Cats
// @Failure 400 {object} models.Cats
// @Router /cats [post]
func (h *CatHandler) CreateCats(c echo.Context) error {
	cats := new(models.Cats)
	err := json.NewDecoder(c.Request().Body).Decode(cats)
	if err != nil {
		return c.JSON(http.StatusBadRequest, new(models.Cats))
	}
	if err = c.Validate(cats); err != nil {
		return c.JSON(http.StatusBadRequest, new(models.Cats))
	}

	cat, err := h.src.CreateCatsServ(*cats)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusCreated, cat)
}

// @Summary GetCat
// @Tags Cats
// @Description get cat by id
// @Accept json
// @Produce json
// @Param id path int true "id"
// @Success 200 {object} models.Cats
// @Failure 400 {object} models.Cats
// @Failure 500 {string} string
// @Router /cats/{id} [get]
func (h *CatHandler) GetCat(c echo.Context) error {
	var catModel partOfCats

	// Достаём ID
	id := c.Param("id")

	// Валидируем
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, new(models.Cats))
	}
	catModel.ID = int32(idInt)
	if err = c.Validate(catModel); err != nil {
		return c.JSON(http.StatusBadRequest, new(models.Cats))
	}

	cat, err := h.src.GetCatServ(id)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, cat)
}

// @Summary UpdateCat
// @Tags Cats
// @Description update cat by id
// @Accept json
// @Produce json
// @Param id path int true "id"
// @Param cats body models.Cats true "cats"
// @Success 200 {object} models.Cats
// @Failure 400 {object} models.Cats
// @Failure 500 {string} string
// @Router /cats/{id} [put]
func (h *CatHandler) UpdateCat(c echo.Context) error {
	cats := new(models.Cats)

	err := json.NewDecoder(c.Request().Body).Decode(&cats)
	if err != nil {
		return c.JSON(http.StatusBadRequest, new(models.Cats))
	}

	// Достаём ID
	id := c.Param("id")

	// Конвертируем id в int
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, new(models.Cats))
	}

	// Присваиваем id модели
	cats.ID = int32(idInt)

	// Проверка валидности модели
	if err = c.Validate(cats); err != nil {
		return c.JSON(http.StatusBadRequest, new(models.Cats))
	}

	cat, err := h.src.UpdateCatServ(id, *cats)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, cat)
}

// @Summary DeleteCat
// @Tags Cats
// @Description delete cat by id
// @Accept json
// @Produce json
// @Param id path int true "id"
// @Success 200 {object} models.Cats
// @Failure 400 {object} models.Cats
// @Failure 500 {string} string
// @Router /cats/{id} [delete]
func (h *CatHandler) DeleteCat(c echo.Context) error {
	var catModel partOfCats

	// Достаём ID
	id := c.Param("id")

	// Валидируем
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, new(models.Cats))
	}
	catModel.ID = int32(idInt)
	if err := c.Validate(catModel); err != nil {
		return c.JSON(http.StatusBadRequest, new(models.Cats))
	}

	cat, err := h.src.DeleteCatServ(id)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, cat)
}
