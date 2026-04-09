package handlers

import (
	"log"
	"net/http"

	"github.com/USACE/go-consequences/consequences"
	"github.com/USACE/go-consequences/geography"
	"github.com/USACE/go-consequences/resultswriters"
	"github.com/USACE/go-consequences/structureprovider"
	"github.com/USACE/go-consequences/structures"
	"github.com/labstack/echo/v5"
	"github.com/usace-nsi/consequences-server/models"
)

type Handler struct {
}
type Status struct {
	State string `json:"state"`
}

func (h *Handler) Version(c *echo.Context) error {
	return c.String(http.StatusOK, "consequences-server:v1.0.0")
}
func (h *Handler) Status(c *echo.Context) error {
	s := Status{State: "RUNNING"}
	return c.JSON(http.StatusOK, s)
}

func (h *Handler) Compute(c *echo.Context) error {
	//
	var data models.ComputeInput
	if err := c.Bind(&data); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}
	hp := models.InitRasDepthJsonProvider(data)
	sp := structureprovider.InitNSISP()
	jrw := resultswriters.InitGeoJsonResultsWriter(c.Response())
	defer jrw.Close()

	bbox, err := hp.HazardBoundary()
	if err != nil {
		return c.JSON(http.StatusExpectationFailed, map[string]string{"error": "Unable to get the raster bounding box"})
	}

	sp.ByBbox(bbox, func(f consequences.Receptor) {
		//ProvideHazard works off of a geography.Location
		s, sok := f.(structures.StructureStochastic)
		if !sok {
			log.Print("nope")
		}
		//leverage srid to pass in the structure name (fdid) to ensure that uniqueness is maintained.
		d, err2 := hp.Hazard(geography.Location{X: f.Location().X, Y: f.Location().Y, SRID: s.Name})
		//compute damages based on hazard being able to provide depth
		if err2 == nil {
			r, err3 := f.Compute(d)
			if err3 == nil {
				jrw.Write(r)
			}
		}
	})
	return nil
}
