package models

import (
	"errors"

	"github.com/USACE/go-consequences/geography"
	"github.com/USACE/go-consequences/hazardproviders"
	"github.com/USACE/go-consequences/hazards"
)

type RasDepthJsonProvider struct {
	data    []DepthEvent
	bbox    geography.BBox
	Process hazardproviders.HazardFunction
}

func InitRasDepthJsonProvider(input ComputeInput) *RasDepthJsonProvider {
	if len(input.HazardData) == 0 {
		return &RasDepthJsonProvider{
			data: input.HazardData,
			bbox: geography.BBox{Bbox: []float64{0, 0, 0, 0}},
		}
	}
	return &RasDepthJsonProvider{
		data: input.HazardData,
		bbox: geography.BBox{
			// Order matches your Contains() logic: [0]=minX, [1]=minY, [2]=maxX, [3]=maxY
			Bbox: input.BoundingBox,
		},
		Process: hazardproviders.DepthHazardFunction(),
	}

}

// Hazard satisfies the HazardProvider interface
func (p *RasDepthJsonProvider) Hazard(l geography.Location) (hazards.HazardEvent, error) {

	for _, d := range p.data {
		var h hazards.HazardEvent
		if d.FD_ID == l.SRID { //use the location SRID to carry the structure fdid to ensure unique matching
			hd := hazards.HazardData{
				Depth: d.Depth,
			}
			h, err := p.Process(hd, h)
			if err != nil {
				return h, err
			}
			return h, nil
		}

	}
	return nil, errors.New("no hazard found at this location")
}

// HazardBoundary satisfies the HazardProvider interface
func (p *RasDepthJsonProvider) HazardBoundary() (geography.BBox, error) {
	return p.bbox, nil
}

// Close satisfies the HazardProvider interface
func (p *RasDepthJsonProvider) Close() {
	// No resources to release for a memory-backed provider
}
