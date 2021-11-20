// Package handler takes parameters from http request and sends it in service
package handler

import (
	"CatsCrud/internal/models"
	"CatsCrud/internal/request"
	"CatsCrud/internal/service"

	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"

	"net/http"
	"strconv"
)

// Handler has an interface of service
type Handler struct {
	src service.Service
}

// NewCatHandler is constructor
func NewCatHandler(srv service.Service) *Handler {
	return &Handler{src: srv}
}

// GetAllCats gets http request, calls func in service and sends http response
// @Summary GetAllCats
// @Tags Cats
// @Description collect all cats in array
// @Produce json
// @Success 200 {array} models.Cats
// @Router /cats [get]
func (h *Handler) GetAllCats(c echo.Context) error {
	everyone, err := h.src.GetAllCatsServ()
	if err != nil {
		log.Error(err)
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, everyone)
}

// CreateCats gets http request, calls func in service and sends http response
// @Summary CreateCats
// @Tags Cats
// @Description create cat
// @Accept json
// @Produce json
// @Param cats body models.Cats true "cats"
// @Success 201 {object} models.Cats
// @Failure 400 {object} models.Cats
// @Router /cats [post]
func (h *Handler) CreateCats(c echo.Context) error {
	cats := new(models.Cats)
	if err := c.Bind(cats); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	if err := c.Validate(cats); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	cat, err := h.src.CreateCatsServ(*cats)
	if err != nil {
		log.Error(err)
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusCreated, cat)
}

// GetCat gets http request, calls func in service and sends http response
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
func (h *Handler) GetCat(c echo.Context) error {
	id := new(request.CatID)
	if err := c.Bind(id); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	if err := c.Validate(id); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	cat, err := h.src.GetCatServ(strconv.Itoa(int(id.ID)))
	if err != nil {
		log.Error(err)
		return c.JSON(http.StatusBadRequest, err)
	}
	return c.JSON(http.StatusOK, cat)
}

// UpdateCat gets http request, calls func in service and sends http response
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
func (h *Handler) UpdateCat(c echo.Context) error {
	cats := new(models.Cats)
	if err := c.Bind(cats); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	if err := c.Validate(cats); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	cat, err := h.src.UpdateCatServ(strconv.Itoa(int(cats.ID)), *cats)
	if err != nil {
		log.Error(err)
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, cat)
}

// DeleteCat gets http request, calls func in service and sends http response
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
func (h *Handler) DeleteCat(c echo.Context) error {
	id := &request.CatID{}
	if err := c.Bind(id); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	if err := c.Validate(id); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	cat, err := h.src.DeleteCatServ(strconv.Itoa(int(id.ID)))
	if err != nil {
		log.Error(err)
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, cat)
}
