package models

import (
	"errors"
	"math"

	"github.com/USACE/go-consequences/geography"
	"github.com/USACE/go-consequences/hazardproviders"
	"github.com/USACE/go-consequences/hazards"
)

type RasDepthJsonProvider struct {
	data    []RasDepths
	bbox    geography.BBox
	Process hazardproviders.HazardFunction
}

func InitRasDepthJsonProvider(depths []RasDepths) *RasDepthJsonProvider {
	if len(depths) == 0 {
		return &RasDepthJsonProvider{
			data: depths,
			bbox: geography.BBox{Bbox: []float64{0, 0, 0, 0}},
		}
	}

	minX, minY := math.MaxFloat64, math.MaxFloat64
	maxX, maxY := -math.MaxFloat64, -math.MaxFloat64

	for _, d := range depths {
		if d.X < minX {
			minX = d.X
		}
		if d.X > maxX {
			maxX = d.X
		}
		if d.Y < minY {
			minY = d.Y
		}
		if d.Y > maxY {
			maxY = d.Y
		}
	}

	return &RasDepthJsonProvider{
		data: depths,
		bbox: geography.BBox{
			// Order matches your Contains() logic: [0]=minX, [1]=minY, [2]=maxX, [3]=maxY
			Bbox: []float64{minX, minY, maxX, maxY},
		},
		Process: hazardproviders.DepthHazardFunction(),
	}

}

// Hazard satisfies the HazardProvider interface
func (p *RasDepthJsonProvider) Hazard(l geography.Location) (hazards.HazardEvent, error) {
	// Optional: Check BBox first to fail fast
	if !p.bbox.Contains(l) {
		return nil, errors.New("location outside of provider bounds")
	}

	for _, d := range p.data {
		// Using a small epsilon for float comparison if necessary,
		// or direct comparison if coordinates are exact matches.
		var h hazards.HazardEvent
		if d.X == l.X && d.Y == l.Y {
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
