package rmu

import (
	"time"
)

type WindSpeedDao interface {
	InsertWindSpeed(id string, value float64, createdAt time.Time) error
	UpdateWindSpeed(id string, value float64, updatedAt time.Time) error
}

type WindSpeed struct {
	id   string
	value float64
}
