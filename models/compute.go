package models

type DepthEvent struct {
	FD_ID string  `db:"fd_id" json:"fd_id"`
	Depth float64 `db:"depth" json:"depth"`
}
type ComputeInput struct {
	BoundingBox     []float64    `db:"bbox" json:"bbox"`
	InventorySource string       `db:"inventory_source" json:"inventory_source"`
	HazardData      []DepthEvent `db:"hazard_data" json:"hazard_data"`
}
