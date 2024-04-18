package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type User struct {
	ID        uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt pgtype.Timestamptz
	Email     string
	Password  string
	Role      int
	CartID    uuid.UUID
	Cart      Cart
	FirstName string
	LastName  string
}
