package entities

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

// Stock - a declared entity.
type Stock struct {
	bun.BaseModel `bun:"table:stock,alias:stock"`

	ID        uuid.UUID `bun:"id,pk,notnull" json:"id" yaml:"id"`
	Name      string    `bun:"name,notnull" json:"name" yaml:"name"`
	Quantity  int64     `bun:"quantity,notnull,nullzero,default:0" json:"quantity" yaml:"quantity"`
	CreatedAt time.Time `bun:",nullzero,notnull,default:current_timestamp" json:"created_at" yaml:"created_at"`
	UpdatedAt time.Time `bun:",nullzero,notnull,default:current_timestamp" json:"updated_at" yaml:"updated_at"`
}

// StockItems a slice of stock entities.
type StockItems []*Stock

// BeforeAppendModel DB hooks that will be executed before a DB query.
func (s *Stock) BeforeAppendModel(_ context.Context, query bun.Query) error {
	switch query.(type) {
	case *bun.InsertQuery:
		if s.ID == uuid.Nil {
			s.ID = uuid.New()
		}

		s.CreatedAt = time.Now()
	case *bun.UpdateQuery:
		s.UpdatedAt = time.Now()
	}
	return nil
}
