package test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
	"stocks-api/module/entities"
)

// TestBeforeAppend asserts that the model's pre-db hooks are working.
func TestBeforeAppend(t *testing.T) {
	createdAt := time.Now()
	updatedAt := time.Now()

	stock := entities.Stock{
		Name:      "test_stock",
		Quantity:  10,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}

	insertQuery := bun.NewInsertQuery(new(bun.DB))

	stock.BeforeAppendModel(context.Background(), insertQuery)

	if stock.ID == uuid.Nil {
		t.Fatal("Expected uuid to be generated, received nil")
	}

	if stock.CreatedAt == createdAt {
		t.Fatalf("Expected timestamp generation diff: %s <> %s", stock.CreatedAt.String(), createdAt.String())
	}

	updateQuery := bun.NewUpdateQuery(new(bun.DB))

	stock.BeforeAppendModel(context.Background(), updateQuery)

	if stock.UpdatedAt == updatedAt {
		t.Fatalf("Expected timestamp generation diff: %s <> %s", stock.UpdatedAt.String(), updatedAt.String())
	}
}
