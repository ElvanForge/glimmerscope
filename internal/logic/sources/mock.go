package sources

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/ElvanForge/glimmerscope/internal/models"
)

type MockSource struct{}

func (s *MockSource) Name() string {
	return "MockMarket"
}

func (s *MockSource) Search(ctx context.Context, query string) ([]models.LorcanaCard, error) {
	// Simulate a slight network delay (Tactical Latency)
	select {
	case <-time.After(100 * time.Millisecond):
	case <-ctx.Done():
		return nil, ctx.Err()
	}

	cards := []models.LorcanaCard{}
	rarities := []string{"Common", "Uncommon", "Rare", "Super Rare", "Legendary", "Enchanted"}

	// Generate 10 mock cards with varying timestamps
	for i := 1; i <= 10; i++ {
		cards = append(cards, models.LorcanaCard{
			ID:        1000 + i,
			Name:      fmt.Sprintf("%s - %s #%d", query, "Glimmer", i),
			Expansion: "The First Chapter",
			Rarity:    rarities[rand.Intn(len(rarities))],
			PriceGuide: models.PriceData{
				Trend: 10.0 + rand.Float64()*50.0,
				Low:   5.0 + rand.Float64()*10.0,
				Avg1:  12.0 + rand.Float64()*40.0,
			},
			// Stagger timestamps so the heap has to actually sort
			UpdatedAt: time.Now().Add(time.Duration(-rand.Intn(60)) * time.Minute),
		})
	}

	return cards, nil
}