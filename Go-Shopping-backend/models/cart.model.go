package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type Cart struct {
	ID        uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt pgtype.Timestamptz `db:"deleted_at"`
	UserID    uuid.UUID
	Products  []CartProduct
}

type CartProduct struct {
	ProductID uuid.UUID
	Quantity  int
}
