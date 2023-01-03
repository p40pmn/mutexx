package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type handler struct {
	svc *Service
}

func (h *handler) Create(c echo.Context) error {
	a, err := h.svc.Create(c.Request().Context())
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	return c.JSON(http.StatusCreated, a)
}
