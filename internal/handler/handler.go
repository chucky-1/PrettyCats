// Package handler takes parameters from http request and sends it in service
package handler

import (
	"CatsCrud/internal/models"
	"CatsCrud/internal/request"
	"CatsCrud/internal/service"
	"github.com/golang-jwt/jwt"

	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"

	"net/http"
	"strconv"
)

// Handler has an interface of service
type Handler struct {
	srv service.Service
}

// NewHandler is constructor
func NewHandler(srv service.Service) *Handler {
	return &Handler{srv: srv}
}

// GetAll gets http request, calls func in service and sends http response
// @Summary GetAllCats
// @Tags Cats
// @Description collect all cats in array
// @Produce json
// @Success 200 {array} models.Cats
// @Router /cats [get]
func (h *Handler) GetAll(c echo.Context) error {
	everyone, err := h.srv.GetAll()
	if err != nil {
		log.Error(err)
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, everyone)
}

// Create gets http request, calls func in service and sends http response
// @Summary CreateCats
// @Tags Cats
// @Description create cat
// @Accept json
// @Produce json
// @Param cats body models.Cats true "cats"
// @Success 201 {object} models.Cats
// @Failure 400 {object} models.Cats
// @Router /cats [post]
func (h *Handler) Create(c echo.Context) error {
	cats := new(models.Cat)
	if err := c.Bind(cats); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	if err := c.Validate(cats); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	cat, err := h.srv.Create(*cats)
	if err != nil {
		log.Error(err)
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusCreated, cat)
}

// Get gets http request, calls func in service and sends http response
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
func (h *Handler) Get(c echo.Context) error {
	id := new(request.CatID)
	if err := c.Bind(id); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	if err := c.Validate(id); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	cat, err := h.srv.Get(strconv.Itoa(int(id.ID)))
	if err != nil {
		log.Error(err)
		return c.JSON(http.StatusBadRequest, err)
	}
	return c.JSON(http.StatusOK, cat)
}

// Update gets http request, calls func in service and sends http response
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
func (h *Handler) Update(c echo.Context) error {
	cats := new(models.Cat)
	if err := c.Bind(cats); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	if err := c.Validate(cats); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	cat, err := h.srv.Update(strconv.Itoa(int(cats.ID)), *cats)
	if err != nil {
		log.Error(err)
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, cat)
}

// Delete gets http request, calls func in service and sends http response
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
func (h *Handler) Delete(c echo.Context) error {
	id := &request.CatID{}
	if err := c.Bind(id); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	if err := c.Validate(id); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	cat, err := h.srv.Delete(strconv.Itoa(int(id.ID)))
	if err != nil {
		log.Error(err)
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, cat)
}

// Restricted is func which check an authorization
// @Summary Restricted
// @Security ApiKeyAuth
// @Description example closed page
// @Produce json
// @Success 200 {string} string
// @Failure 400 {object} string
// @Router /restrict [get]
func (h *AuthHandler) Restricted(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*service.JwtCustomClaims)
	name := claims.Name
	id := claims.ID
	idStr := strconv.Itoa(id)
	return c.String(http.StatusOK, "Welcome "+name+idStr)
}

