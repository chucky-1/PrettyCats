package handler

import (
	"CatsCrud/internal/models"
	"CatsCrud/internal/request"
	"CatsCrud/internal/service"
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
	if err := c.Bind(cats); err != nil {
		return c.JSON(http.StatusBadRequest, new(models.Cats))
	}
	if err := c.Validate(cats); err != nil {
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
	id := new(request.RequestCatId)
	if err := c.Bind(id); err != nil {
		return err
	}
	if err := c.Validate(id); err != nil {
		return c.JSON(http.StatusBadRequest, new(models.Cats))
	}

	cat, err := h.src.GetCatServ(strconv.Itoa(int(id.ID)))
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
	if err := c.Bind(cats); err != err {
		return c.JSON(http.StatusBadRequest, new(models.Cats))
	}
	if err := c.Validate(cats); err != nil {
		return c.JSON(http.StatusBadRequest, new(models.Cats))
	}

	cat, err := h.src.UpdateCatServ(strconv.Itoa(int(cats.ID)), *cats)
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
	id := new(request.RequestCatId)
	if err := c.Bind(id); err != nil {
		return c.JSON(http.StatusBadRequest, new(models.Cats))
	}
	if err := c.Validate(id); err != nil {
		return c.JSON(http.StatusBadRequest, new(models.Cats))
	}

	cat, err := h.src.DeleteCatServ(strconv.Itoa(int(id.ID)))
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, cat)
}
