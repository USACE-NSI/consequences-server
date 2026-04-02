package models

type RasDepths struct {
	FD_ID string  `db:"fd_id" json:"fd_id"`
	X     float64 `db:"x" json:"x"`
	Y     float64 `db:"y" json:"y"`
	Depth float64 `db:"depth" json:"depth"`
}
