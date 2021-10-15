package handler

import (
	"CatsCrud/service"
	"github.com/labstack/echo/v4"
	"net/http"
)

type CatHandler struct {
	src *service.CatService
}

func NewCatHandler(srv *service.CatService) *CatHandler {
	return &CatHandler{src: srv}
}

func(h *CatHandler) GetAllCats(c echo.Context) error {
	allcats, err :=  h.src.GetAllCatsServ()
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, allcats)
}

//func (h *CatHandler) CreateCats(c echo.Context) error {
//
//	// Получаем параметры json. Здесь, потому что repository не имеет доступа к echo
//	jsonMap := make(map[string]interface{})
//	err := json.NewDecoder(c.Request().Body).Decode(&jsonMap)
//	if err != nil {
//		return err
//	}
//
//	cat, err := h.src.CreateCatsServ(jsonMap)
//	if err != nil {
//		return err
//	}
//	return c.JSON(http.StatusCreated, cat)
//}
//
//func (h *CatHandler) GetCat(c echo.Context) error {
//	// Достаём ID
//	id := c.Param("id")
//
//	cat, err := h.src.GetCatServ(id)
//	if err != nil {
//		return err
//	}
//	return c.JSON(http.StatusOK, cat)
//}
//
//func (h *CatHandler) UpdateCat(c echo.Context) error {
//	// Достаём ID
//	id := c.Param("id")
//
//	jsonMap := make(map[string]interface{})
//	err := json.NewDecoder(c.Request().Body).Decode(&jsonMap)
//	if err != nil {
//		return err
//	}
//
//	cat, err := h.src.UpdateCatServ(id, jsonMap)
//	if err != nil {
//		return err
//	}
//	return c.JSON(http.StatusOK, cat)
//}
//
//func (h *CatHandler) DeleteCat(c echo.Context) error {
//	// Достаём ID
//	id := c.Param("id")
//
//	cat, err := h.src.DeleteCatServ(id)
//	if err != nil {
//		return err
//	}
//	return c.JSON(http.StatusOK, cat)
//}
//
