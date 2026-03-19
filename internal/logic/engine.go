package logic

import (
	"container/heap"
	"context"
	"sync"
	"time"

	"github.com/ElvanForge/glimmerscope/internal/models"
)

type Card = models.LorcanaCard

type Source interface {
	Search(ctx context.Context, query string) ([]Card, error)
	Name() string
}

type resultsHeap []Card

func (h resultsHeap) Len() int           { return len(h) }
func (h resultsHeap) Less(i, j int) bool { return h[i].UpdatedAt.Before(h[j].UpdatedAt) }
func (h resultsHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *resultsHeap) Push(x interface{}) { *h = append(*h, x.(Card)) }
func (h *resultsHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

type Engine struct {
	Sources []Source
}

func (e *Engine) Collect(ctx context.Context, query string) []Card {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	h := &resultsHeap{}
	heap.Init(h)
	resultsChan := make(chan []Card, len(e.Sources))
	var wg sync.WaitGroup

	for _, s := range e.Sources {
		wg.Add(1)
		go func(src Source) {
			defer wg.Done()
			cards, err := src.Search(ctx, query)
			if err != nil {
				return
			}
			resultsChan <- cards
		}(s)
	}

	go func() { wg.Wait(); close(resultsChan) }()

Loop:
	for {
		select {
		case cards, ok := <-resultsChan:
			if !ok { break Loop }
			for _, c := range cards {
				if h.Len() < 20 {
					heap.Push(h, c)
				} else if c.UpdatedAt.After((*h)[0].UpdatedAt) {
					heap.Pop(h)
					heap.Push(h, c)
				}
			}
		case <-ctx.Done(): break Loop
		}
	}

	final := make([]Card, h.Len())
	for i := h.Len() - 1; i >= 0; i-- { final[i] = heap.Pop(h).(Card) }
	return final
}