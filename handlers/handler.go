package handlers

import (
	"net/http"

	"github.com/USACE/go-consequences/compute"
	"github.com/USACE/go-consequences/resultswriters"
	"github.com/USACE/go-consequences/structureprovider"
	"github.com/labstack/echo/v5"
	"github.com/usace-nsi/consequences-server/models"
)

type Handler struct {
}

func (h *Handler) Version(c *echo.Context) error {
	return c.String(http.StatusOK, "consequences-server:v0.0.1")
}

func (h *Handler) Compute(c *echo.Context) error {
	//
	var data []models.RasDepths
	if err := c.Bind(&data); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}
	hp := models.InitRasDepthJsonProvider(data)
	sp := structureprovider.InitNSISP()
	jrw := resultswriters.InitGeoJsonResultsWriter(c.Response())
	defer jrw.Close()
	compute.StreamAbstract(hp, sp, jrw)

	return nil
}
