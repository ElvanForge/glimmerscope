package models

import "time"

type LorcanaCard struct {
    ID          int       `json:"idProduct"`
    Name        string    `json:"enName"`
    Expansion   string    `json:"expansionName"`
    Rarity      string    `json:"rarity"`
    PriceGuide  PriceData `json:"priceGuide"`
    UpdatedAt   time.Time `json:"updatedAt"`
}

type PriceData struct {
    Trend float64 `json:"TREND"`
    Low   float64 `json:"LOW"`
    Avg1  float64 `json:"AVG1"`
}