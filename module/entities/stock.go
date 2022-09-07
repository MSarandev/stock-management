package entities

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

// Stock - a declared entity.
type Stock struct {
	bun.BaseModel `bun:"table:stocks,alias:stocks"`

	ID        uuid.UUID `bun:"id,pk,notnull" json:"id" yaml:"id"`
	Name      string    `bun:"name,notnull" json:"name" yaml:"name"`
	Quantity  int64     `bun:"quantity,notnull,nullzero,default:0" json:"quantity" yaml:"quantity"`
	CreatedAt time.Time `bun:",nullzero,notnull,default:current_timestamp" json:"created_at" yaml:"created_at"`
	UpdatedAt time.Time `bun:",nullzero,notnull,default:current_timestamp" json:"updated_at" yaml:"updated_at"`
}
