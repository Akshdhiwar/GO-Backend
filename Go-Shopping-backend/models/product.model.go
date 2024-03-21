package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type Product struct {
	ID          uuid.UUID
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   pgtype.Timestamptz `db:"deleted_at"`
	Title       string
	Price       float64
	Description string
	Category    string
	Image       string
	Rating      float32
	Count       int
}
